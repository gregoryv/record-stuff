/* Copyright 2013 Chris Wilson

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

window.AudioContext = window.AudioContext || window.webkitAudioContext;

var audioContext = new AudioContext();
var audioInput = null,
    realAudioInput = null,
    inputPoint = null,
    audioRecorder = null;
var rafID = null;
var analyserContext = null;
var canvasWidth, canvasHeight;



// Make the function wait until the connection is made...
function waitForSocketConnection(socket, callback){
    setTimeout(
        function () {
            if (socket.readyState === 1) {
                console.log("Connection is made")
                if(callback != null){
                    callback();
                }
                return;

            } else {
                console.log("wait for connection...")
                waitForSocketConnection(socket, callback);
            }

        }, 5); // wait 5 milisecond for the connection...
}


function postSound(blob) {
    var fd = new FormData();
    var timestamp = new Date();
    var filename = timestamp.getTime() + ".wav";
    fd.append('filename', filename);
    fd.append('soundBlob', blob);
    $.ajax({
	type: 'POST',
	url: '/upload',
	data: fd,
	processData: false,
	contentType: false
    }).done(function(data) {
	console.log(data);
	var href = "/recordings/"+ filename;
	newLi(href, timestamp);
    });
}

function newLi(href, name) {
    $("#recordings").append("<li><a href=\"" + href + "\">" + name + "</a></li>");
}

var ws;
var chunkerId;

function toggleRecording( e ) {
    var useWebsocket = true;
    if (e.classList.contains("recording")) {
        e.classList.remove("recording");
	
	if(useWebsocket) {
	    audioRecorder.stop();
	    clearInterval(chunkerId);
	    console.log("Closing socket connection");
	    ws.close();
	} else {
            audioRecorder.stop();
	    audioRecorder.exportWAV( postSound );
	}
    } else {
        // start recording
        if (!audioRecorder)
            return;
        e.classList.add("recording");
	if(useWebsocket) {
	    // sound is send in chunks 
	    ws = new WebSocket("ws://localhost:8080/record");
	    ws.binaryType = "blob"
	    waitForSocketConnection(ws, function(){
		chunkerId = setInterval( function(){
		    // Issue here, cannot export unless recorder is stopped, not good!
		    audioRecorder.stop();
		    audioRecorder.exportWAV(function(chunk) {
			audioRecorder.clear();						
			audioRecorder.record();
			console.log("Send chunk");			
			ws.send(chunk);
		    });

		}, 1000);
	    });
	} else {
            audioRecorder.clear();
            audioRecorder.record();
	}
    }
}

// Called when microphone is activated
function gotStream(stream) {
    inputPoint = audioContext.createGain();

    // Create an AudioNode from the stream.
    realAudioInput = audioContext.createMediaStreamSource(stream);
    audioInput = realAudioInput;
    audioInput.connect(inputPoint);

    analyserNode = audioContext.createAnalyser();
    analyserNode.fftSize = 2048;
    inputPoint.connect( analyserNode );

    audioRecorder = new Recorder( inputPoint );

    zeroGain = audioContext.createGain();
    zeroGain.gain.value = 0.0;
    inputPoint.connect( zeroGain );
    zeroGain.connect( audioContext.destination );
}

function initAudio() {
        if (!navigator.getUserMedia)
            navigator.getUserMedia = navigator.webkitGetUserMedia || navigator.mediaDevices.getUserMedia;
        if (!navigator.cancelAnimationFrame)
            navigator.cancelAnimationFrame = navigator.webkitCancelAnimationFrame || navigator.mozCancelAnimationFrame;
        if (!navigator.requestAnimationFrame)
            navigator.requestAnimationFrame = navigator.webkitRequestAnimationFrame || navigator.mozRequestAnimationFrame;

    navigator.getUserMedia(
        {
            "audio": {
                "mandatory": {
                    "googEchoCancellation": "false",
                    "googAutoGainControl": "false",
                    "googNoiseSuppression": "false",
                    "googHighpassFilter": "false"
                },
                "optional": []
            },
        }, gotStream, function(e) {
            alert('Error getting audio');
            console.log(e);
        });
}

function loadRecordings() {
 $(function() {
     $.getJSON("/recordings/", function(data) {
	 $(data).each(function(i, el) {
	     var parts = el.name.split(".")
	     newLi(el.href, new Date(+parts[0]));
	 });
     });
 });
}

window.addEventListener('load', initAudio );
