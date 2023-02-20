package data

import (
	v1 "chatbot/api/helloworld/v1"
	"chatbot/pkg/define"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type openai struct {
	token string
	proxy string
	code  string
	conf  config
}
type config struct {
	Model            string           `json:"model"`
	Prompt           *string          `json:"prompt,omitempty"`
	Suffix           *string          `json:"suffix,omitempty"`
	MaxTokens        *int64           `json:"max_tokens,omitempty"`
	Temperature      *float64         `json:"temperature,omitempty"`
	TopP             *float32         `json:"top_p,omitempty"`
	N                *int64           `json:"n,omitempty"`
	Stream           *bool            `json:"stream,omitempty"`
	Logprobs         *int64           `json:"logprobs,omitempty"`
	Echo             *bool            `json:"echo,omitempty"`
	Stop             []string         `json:"stop,omitempty"`
	PresencePenalty  *float32         `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float32         `json:"frequency_penalty,omitempty"`
	BestOf           *int64           `json:"best_of,omitempty"`
	LogitBias        map[string]int64 `json:"logit_bias,omitempty"`
	User             *string          `json:"user,omitempty"`
}

func (o *openai) client() *gogpt.Client {
	c := gogpt.NewClient(o.token)
	if o.proxy != "" {
		pu, _ := url.Parse(o.proxy)
		c.HTTPClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(pu),
		}
	}
	return c
}
func (o *openai) request(q string) gogpt.CompletionRequest {
	req := gogpt.CompletionRequest{}
	tmp, _ := json.Marshal(o.conf)
	_ = json.Unmarshal(tmp, &req)
	req.Prompt = q
	return req
}

const ErrCodeInvalid = "bot code invalid"

func (o *openai) send(ctx context.Context, q, code string) (message define.Message, err error) {
	t1 := time.Now()
	text := ErrCodeInvalid
	resp := gogpt.CompletionResponse{}
	if o.checkCode(code) {
		resp, err = o.client().CreateCompletion(ctx, o.request(q))
		if err == nil {
			text = resp.Choices[0].Text
		}
	}
	t2 := time.Since(t1).Milliseconds()
	message = define.Message{
		Id:   define.GenID(),
		Time: t2,
		Eof:  true,
		Err:  err,
	}
	if err != nil {
		return
	}
	message.Msg = text
	return
}

func (o *openai) testData(q string, ch chan define.Message) {
	defer close(ch)
	str := "床前明月光，疑是地上霜。\n举头望明月，低头思故乡。hello,my god."
	t1 := define.GenID()
	for _, v := range str {
		m := define.Message{
			Id:  t1,
			Msg: string(v),
		}
		ch <- m
	}
	ch <- define.Message{
		Id:  t1,
		Msg: "",
		Eof: true,
	}
}
func (o *openai) checkCode(code string) bool {
	if o.code != "" && o.code != strings.TrimSpace(code) {
		return false
	}
	return true
}

func (o *openai) stream(realQ define.Send, ch chan define.Message) {
	var (
		response gogpt.CompletionResponse
		stream   *gogpt.CompletionStream
		end      define.Message
		err      error
		id       int64
		costTime int64
		cancel   context.CancelFunc
	)
	t1 := time.Now()
	id = define.GenID()
	defer func() {
		costTime = time.Since(t1).Milliseconds()
		end.Id = id
		end.Time = costTime
		end.Eof = true
		end.Err = err
		ch <- end
		close(ch)
	}()
	if !o.checkCode(realQ.BotCode) {
		err = v1.ErrorCodeInvalid("bot code错误")
		return
	}
	q := strings.TrimSpace(realQ.Msg)
	ctx := context.Background()
	ctx, cancel = context.WithTimeout(ctx, 5*60*time.Second)
	defer cancel()
	if q == "" {
		err = v1.ErrorQuestionNull("问题为空")
		return
	}
	stream, err = o.client().CreateCompletionStream(ctx, o.request(q))
	if err != nil {
		err = v1.ErrorGreeterUnspecified("openAi:%s", err.Error())
		return
	}
	defer stream.Close()
	for {
		response, err = stream.Recv()
		if errors.Is(err, io.EOF) {
			err = nil
			return
		}
		if err != nil {
			err = v1.ErrorGreeterUnspecified("openAi:%s", err.Error())
			return
		}
		if len(response.Choices) > 0 {
			ch <- define.Message{
				Id:   id,
				Msg:  response.Choices[0].Text,
				Time: time.Since(t1).Milliseconds(),
			}
		}
	}

}
