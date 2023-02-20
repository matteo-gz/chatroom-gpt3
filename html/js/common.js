const HTMLEncode = function (html) {
    let temp = document.createElement("div");
    (temp.textContent != null) ? (temp.textContent = html) : (temp.innerText = html);
    let output = temp.innerHTML;
    temp = null;
    return output;
}
const nowt = function () {
    let date = new Date();
    let hours = date.getHours();
    let minutes = date.getMinutes();
    let seconds = date.getSeconds();
    return hours + ":" + minutes + ":" + seconds
}
let writeCookie = function (name, value) {
    let days = 30;
    let expires = new Date();
    expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
    document.cookie = name + "=" + escape(value) + ";expires=" + expires.toGMTString();
}
let getCookie = function (cname) {
    let name = cname + "=";
    let ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i].trim();
        if (c.indexOf(name) === 0) return c.substring(name.length, c.length);
    }
    return "";
}
const synth = window.speechSynthesis;

const textToSpeech = (text, index) => {
    let voices = synth.getVoices();
    let i = 0
    let femaleVoice = voices.find(voice => {
        if (index == i) {
            return voice.name
        }
        i++
    });
    let utterThis = new SpeechSynthesisUtterance(text);
    utterThis.rate = 1.0;
    utterThis.voice = femaleVoice;
    synth.cancel()
    synth.speak(utterThis);
};
let VoiceList = function () {
    if (typeof speechSynthesis === 'undefined') {
        console.log("no")
        return;
    }
    let voices = speechSynthesis.getVoices();
    let l = []
    for (let i = 0; i < voices.length; i++) {
        let option = document.createElement('option');
        option.textContent = `${voices[i].name} (${voices[i].lang})`;

        if (voices[i].default) {
            option.textContent += ' — DEFAULT';
        }
        option.setAttribute('data-lang', voices[i].lang);
        option.setAttribute('data-name', voices[i].name);
        l.push(option)
    }
    return l;
}

const SpeechRecognition = window.SpeechRecognition || webkitSpeechRecognition;
const SpeechGrammarList = window.SpeechGrammarList || webkitSpeechGrammarList;
const SpeechRecognitionEvent = window.SpeechRecognitionEvent || webkitSpeechRecognitionEvent;
const recognition = new SpeechRecognition();
recognition.interimResults = false;
recognition.continuous = true;