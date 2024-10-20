document.addEventListener('paste', async (e: ClipboardEvent) => {
    e.preventDefault();
    if (!e.clipboardData) return
    if (e.clipboardData.files.length > 0) {
        const data = e.clipboardData.files;
        postIt("/upload", data[0]);
    } else {
        const text = e.clipboardData.getData("text/plain");
        postIt("/upload", text);
    }
  });

window.addEventListener('load', async (e) => {
    const urlParams = new URLSearchParams(window.location.search);
    const id = urlParams.get("id");
    if (id != undefined) {
        showContent(id);
    }
  });


async function postIt(url: string, data: any) {
    const response = await fetch(url, {method: "POST", body: data})
    if (response.ok) {
        const text = await response.text();
        const ctype = response.headers.get("Content-Type");
        if (ctype != "redirect") {
            setContent(`<a href="/?id=${text}">Click </a>`);
            return;
        }
        setContent(text, ctype);
    }
}

async function getIt(url: string) {
    const response = await fetch(url);
    const blob = await response.blob()
    setContent(blob, blob.type)
}

async function showContent(id: string) {
    const data = await getIt(`/get?id=${id}`);
    if (data != undefined) setContent(data);
}

async function showFullscreen() {
    const overlay = document.getElementById('fullscreenOverlay');
    const fullscreenImg = document.getElementById('fullscreenImage');
    const image = document.getElementById('image');
    const fsImage = document.getElementById('fullscreenImage');
    if (!overlay || ! fullscreenImg || !image || !fsImage) return;
    overlay.style.display = 'flex';
    const src = image.getAttribute('src') || "";
    fsImage.setAttribute('src', src)
    overlay.onclick = function() {
        overlay.style.display = 'none';
    };
}

async function setContent(content: any, ctype?: string | null) {
    console.log("setting content", ctype);
    console.log(content)
    const share = document.getElementById("share");
    const contentBox = document.getElementById("content");
    const footnote = document.getElementById("footnote");
    if (!share || !contentBox) return;
    share.hidden = true;
    contentBox.hidden = false;
    switch(ctype) {
        case "image/png":
            var blob = content as Blob;
            const reader = new FileReader();
            reader.onloadend = () => {
                const bstring = reader.result;
                contentBox.innerHTML = `<img id="image" src="${bstring}" onClick="showFullscreen()">`;
            }
            reader.readAsDataURL(blob);
            break;
        case "redirect":
            contentBox.innerHTML = `<a href="/r/${content}">Click</a>`;
            break;
        default:
            try {
                var blob = content as Blob;
                contentBox.innerHTML = await blob.text();
            } catch (e) {
                contentBox.innerHTML = content;
                if (!footnote) return;
                footnote.innerHTML = `Content-Type not recogized: ${ctype}`
            }
    }
}
