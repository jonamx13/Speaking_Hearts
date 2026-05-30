let socket = null;
const displayArea = document.getElementById('display-area');
const langSelect = document.getElementById('lang-select');
const connectBtn = document.getElementById('connect-btn');

function connect() {
    if (socket) {
        socket.close();
    }

    const lang = langSelect.value;
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?lang=${lang}`;

    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
        displayArea.innerText = `Connected! Waiting for subtitles in [${lang}]...`;
        connectBtn.innerText = 'Reconnect';
    };

    socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            console.log('Raw JSON from Go:', data);
            const selectedLang = langSelect.value;
            
            // Determine which text to display: specific translation or original text
            let textToDisplay = data.original_text;
            if (data.translations && data.translations[selectedLang]) {
                textToDisplay = data.translations[selectedLang].text;
            }

            displayArea.innerHTML = `
                <div>${textToDisplay}</div>
                <div class="timestamp">${data.timestamp} - Speaker: ${data.speaker_id}</div>
            `;
        } catch (e) {
            console.error('Error parsing message:', e);
        }
    };

    socket.onclose = () => {
        displayArea.innerText = 'Disconnected. Please reconnect.';
        connectBtn.innerText = 'Connect';
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        displayArea.innerText = 'Connection error.';
    };
}

connectBtn.addEventListener('click', connect);
