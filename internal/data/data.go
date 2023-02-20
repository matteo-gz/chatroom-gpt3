package data

import (
	"chatbot/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	*openai
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	d := &Data{
		openai: &openai{
			token: c.Openai.Token,
			proxy: c.Openai.Proxy,
			code:  c.Openai.Code,
			conf: config{
				Model:            c.Openai.Config.Model,
				Prompt:           c.Openai.Config.Prompt,
				Suffix:           c.Openai.Config.Suffix,
				MaxTokens:        c.Openai.Config.MaxTokens,
				Temperature:      c.Openai.Config.Temperature,
				TopP:             c.Openai.Config.TopP,
				N:                c.Openai.Config.N,
				Stream:           c.Openai.Config.Stream,
				Logprobs:         c.Openai.Config.Logprobs,
				Echo:             c.Openai.Config.Echo,
				Stop:             c.Openai.Config.Stop,
				PresencePenalty:  c.Openai.Config.PresencePenalty,
				FrequencyPenalty: c.Openai.Config.FrequencyPenalty,
				BestOf:           c.Openai.Config.BestOf,
				LogitBias:        c.Openai.Config.LogitBias,
				User:             c.Openai.Config.User,
			},
		},
	}
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return d, cleanup, nil
}
