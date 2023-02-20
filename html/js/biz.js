let conn;
let wsState;


let newWs = function () {
    wsState = true
    conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.onclose = function (evt) {
        wsState = false
        tip("webSocket", "Connection closed.")
    };
    conn.onopen = function () {
        console.log("open")
        SendWs(2, "ping")
    }
    conn.onmessage = function (evt) {
        try {
            let arr1 = evt.data.split('\n')
            for (let i = 0; i < arr1.length; i++) {
                scroll1(function () {
                    botSay(arr1[i])
                })
            }
        } catch (err) {
            console.log(err)
        }
    };
}
const toastLiveExample = document.getElementById('liveToast')
const toast = new bootstrap.Toast(toastLiveExample)
let tip = function (title, msg) {
    toastLiveExample.querySelector(".toast-body").innerText = msg
    toastLiveExample.querySelector(".me-auto").innerText = title
    toast.show()
}
window.onload = function () {

    if (window["WebSocket"]) {
        newWs()
        setInterval(function () {
            if (wsState === false) {
                newWs()
            }
        }, 5000)
    } else {
        tip("webSocket", "Your browser does not support WebSockets.")
    }
    populateVoiceList()
    if (getCookie("enter") == 1) {
        document.getElementById('flexSwitchCheckDefault').checked = true
    }
    if (getCookie("speech") == 1) {
        document.getElementById('speech1').checked = true
    }
    if (getCookie("sw2") == 1) {
        document.getElementById('sw2').checked = true
    } else {
        document.getElementById("voice1").hidden = true
    }
    if (getCookie("sw3") == 1) {
        document.getElementById('sw3').checked = true
    } else {
        document.getElementById("submit-btn").hidden = true
    }
    if (getCookie("ask_bot") == 1) {
        document.getElementById('askbot1').checked = true
    }
    if (getCookie('voice') != "") {
        document.getElementById('select1').selectedIndex = getCookie('voice');
    }
    if (getCookie('bot_code') != "") {
        document.getElementById('botcode1').value = getCookie('bot_code');
    }
    if (getCookie('your_name') != "") {
        document.getElementById('yourName1').value = getCookie('your_name');
    } else {
        let nowN = randomString(6)
        writeCookie('your_name', nowN)
        document.getElementById('yourName1').value = nowN
    }
}
let randomString = function (e) {
    e = e || 32;
    let t = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678",
        a = t.length,
        n = "";
    for (let i = 0; i < e; i++) n += t.charAt(Math.floor(Math.random() * a));
    return n
}
let populateVoiceList = function () {
    VoiceList().forEach(element => {
        document.getElementById("select1").appendChild(element);
    })
}

let streamSend = function () {
    let prompt = document.getElementById('form1');
    prompt.value = prompt.value.trim()
    if (!prompt.value) {
        return false;
    }
    if (!conn) {
        return false;
    }
    SendWs(1, prompt.value)
    scroll1(function () {
        appendMsg(getCookie("your_name"), 1, prompt.value)
    })
    prompt.value = ""
}
let SendWs = function (type, value) {
    let jso = {
        "msg": value,
        "types": type,
        "bot_code": getCookie("bot_code"),
        "your_name": getCookie("your_name"),
    }
    conn.send(JSON.stringify(jso));
}
document.getElementById('submit-btn2').addEventListener('click', async function () {
    if (getCookie("ask_bot") == 1) {
        appendAt()
    }
    streamSend()
});
let appendAt = function () {
    let prompt = document.getElementById('form1');
    prompt.value = prompt.value.trim()
    if (!prompt.value) {
        return false;
    }
    prompt.value = "@bot " + prompt.value
}
document.getElementById('submit-btn3').addEventListener('click', async function () {
    appendAt()
    streamSend()
});
document.getElementById('submit-btn').addEventListener('click', async function () {
    let prompt = document.getElementById('form1');
    prompt.value = prompt.value.trim()
    if (!prompt.value) {
        return false;
    }
    scroll1(function () {
        appendMsg(getCookie("your_name"), 1, prompt.value)
        send(prompt.value)
    })
    prompt.value = ""
})


let l1 = document.getElementById('list1');
const scroll1 = function (fn) {
    let log = document.getElementById('f1');
    let doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    fn()
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}
const removeWsId = function () {

}
let wsId
const appendWsId = function (id) {
    wsId = id
}
const botSay = function (reply) {
    let msg = JSON.parse(reply)
    let p2 = document.getElementById(msg.id);
    let s2 = msg.msg
    let havErr = false
    if (msg.err != null) {
        s2 = JSON.stringify(msg.err)
        havErr = true
    }
    if (s2 == "") {
        return
    }
    let costT = Math.floor(msg.time / 1000)
    let myName = getCookie('your_name')
    let nowName = "bot"
    if (msg.types == 1) {
        if (wsId == msg.from) {
            return;
        }
        nowName = msg.your_name + "(" + msg.from + ")"
    } else if (msg.types == 2) {
        appendWsId(msg.msg)
        return;
    } else if (msg.types == 0) {
        // sys
        if (havErr == true) {
            if (wsId != msg.from) {
                return;
            } else {
                tip("error", msg.err.message)
                return;
            }
        }
    } else {

    }

    if (p2) {
        let realMsg = p2.querySelector('.msg-content').querySelector(".msg-reply").querySelector("pre")
        realMsg.innerText += s2
        if (costT > 0) {
            p2.querySelector('.cost').innerText = `[${costT}]`
        }

        if (msg.eof == true) {
            speechText(realMsg.innerText)
        }
    } else {
        let p1 = appendMsg(nowName, 2, s2)
        p1.id = msg.id
        if (costT > 0) {
            p1.querySelector('.cost').innerText = `[${costT}]`
        }
        if (msg.eof == true) {
            speechText(s2)
        }
    }

}
document.getElementById('form1').addEventListener('keydown', (e) => {
    if (e.shiftKey && e.which === 13 && getCookie("enter") == 1) {
        if (getCookie("ask_bot") == 1) {
            appendAt()
        }
        streamSend()
    }
});
document.getElementById('select1').addEventListener('change', function () {
    // let v = this.options[this.selectedIndex].getAttribute('data-name')
    writeCookie("voice", this.selectedIndex)
})
document.getElementById('flexSwitchCheckDefault').addEventListener('change', function () {
    let v = 0
    if (this.checked) {
        v = 1
    }
    writeCookie("enter", v)
});
document.getElementById('yourName1').addEventListener("change", function () {
    this.value = this.value.trim()
    if (this.value == "") {
        this.value = randomString(6)
    }
    writeCookie("your_name", this.value)
})
document.getElementById('botcode1').addEventListener("change", function () {
    writeCookie("bot_code", this.value)
})
document.getElementById("speech1").addEventListener('change', function () {
    let v = 0
    if (this.checked) {
        v = 1
    }
    writeCookie("speech", v)
})
document.getElementById("sw2").addEventListener('change', function () {
    let v = 0
    if (this.checked) {
        v = 1
        document.getElementById("voice1").hidden = false
    } else {
        document.getElementById("voice1").hidden = true
    }
    writeCookie("sw2", v)
})
document.getElementById("sw3").addEventListener('change', function () {
    let v = 0
    if (this.checked) {
        v = 1
        document.getElementById("submit-btn").hidden = false
    } else {
        document.getElementById("submit-btn").hidden = true
    }
    writeCookie("sw3", v)
})
document.getElementById("askbot1").addEventListener('change', function () {
    let v = 0
    if (this.checked) {
        v = 1
    }
    writeCookie("ask_bot", v)
})


const appendMsg = function (user, userT, reply) {
    let p1 = document.createElement('div');
    let nowt1 = nowt()
    let msg = HTMLEncode(reply)

    let userTpl = `
  <div class="msg-wrap">
                    <div class="d-flex justify-content-center msg-time">${nowt1}</div>
                    <div class="msg-box msg-user">
                        <div class="d-flex justify-content-end">
                            <span>${user}</span>
                        </div>

                        <div class="shadow-sm msg-content">
                            <div class="d-flex justify-content-end">
                                <button type="button" class="btn btn-outline-secondary btn-sm" onclick="copy(this)">
                                    <i class="bi-clipboard"></i>
                                </button>
                            </div>
                            <div class="msg-reply text-break">
                            <pre>${msg}</pre>
                            </div>
                        </div>
                    </div>
                </div>
`
    let otherTpl = `
 <div class="msg-wrap">
                    <div class="d-flex justify-content-center msg-time">${nowt1}<span class="cost"></span></div>
                    <div class="msg-box msg-other">
                        <div class="d-flex justify-content-start">
                            <span>${user}</span>
                        </div>

                        <div class="shadow-sm msg-content">
                            <div class="d-flex justify-content-end">
                                <button type="button" class="btn btn-outline-secondary btn-sm" onclick="copy(this)">
                                    <i class="bi-clipboard"></i>
                                </button>
                            </div>
                            <div class="msg-reply text-break">
                            <pre>${msg}</pre></div>
                        </div>
                    </div>
                </div>
`
    if (userT == 1) {
        p1.innerHTML = userTpl;
    } else if (userT == 2) {
        p1.innerHTML = otherTpl;
    }

    p1.className = "msg-box"
    l1.appendChild(p1);
    return p1
}
const speechText = function (reply) {
    if (getCookie("speech") == 1) {
        let def = 0
        if (getCookie("voice") != "") {
            def = getCookie("voice")
        }
        textToSpeech(reply, def);
    }
}


let api_lock = false
const send = function (prompt) {
    if (api_lock === true) {
        return
    }
    api_lock = true
    let myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    let raw = JSON.stringify({
        "q": prompt,
        "code": getCookie("bot_code")
    });
    let requestOptions = {
        method: 'POST', headers: myHeaders, body: raw, redirect: 'follow'
    };
    fetch("http://" + document.location.host + "/question", requestOptions)
        .then(response => {
            api_lock = false
            if (response.status !== 200) {
                throw response.json()
            }
            return response.json()
        })
        .then(result => {
            scroll1(function () {
                botSay(result.res)
            })
        })
        .catch(error => {
            if (error.then != null) {
                error.then(res => {
                    tip("error", res.message)
                })
            } else {
                tip("error", error)
            }
        });
}
let tmpA = ""
let tmpB = ""
document.getElementById("voice1").addEventListener('click', function () {
    if (this.querySelector("i").className === "bi-mic-fill") {
        tmpA = this.querySelector("i").className
        tmpB = document.getElementById("label2").innerText
        this.querySelector("i").className = "bi-stop-circle"
        document.getElementById("label2").innerHTML = "<i class=\"bi-soundwave\"></i>"
        recognition.start();
    } else {
        this.querySelector("i").className = tmpA
        recognition.stop();
        document.getElementById("form1").value += document.getElementById("tmpV2").innerText
        document.getElementById("tmpV2").innerText = ""
        document.getElementById("label2").innerHTML = tmpB
    }

})
recognition.addEventListener('result', (event) => {
    // https://www.ydisp.cn/developer/128419.html
    const results = Array.from(event.results);
    const transcript = results
        .map((result) => result[0])
        .map((result) => result.transcript)
        .join('');
    document.getElementById("tmpV2").innerText = transcript;
});
let copy = function (it) {
    try {
        let t = it.parentElement.parentElement.querySelector(".msg-reply").innerText
        // navigator.clipboard.writeText(t);
        copy2(it, t)
    } catch (err) {
        console.log(err)
    }
}

let copy2 = function (it, e) {
    let text = e;
    let inputElement = document.createElement("input");
    inputElement.value = text;
    document.body.appendChild(inputElement);
    inputElement.select(); //选中文本
    document.execCommand("copy"); //执行浏览器复制命令
    inputElement.remove();
    it.querySelector("i").className = "bi-clipboard-check-fill"
    setTimeout(function () {
        it.querySelector("i").className = "bi-clipboard"
    }, 1000)
}
