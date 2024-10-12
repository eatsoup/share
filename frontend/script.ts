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

async function setContent(content: any, ctype?: string | null) {
    console.log("setting content", ctype);
    console.log(content)
    const element = document.getElementById("content");
    if (!element) return;
    switch(ctype) {
        case "image/png":
            var blob = content as Blob;
            const reader = new FileReader();
            reader.onloadend = () => {
                const bstring = reader.result;
                element.innerHTML = `<img src="${bstring}">`;
            }
            reader.readAsDataURL(blob);
            break;
        case "redirect":
            element.innerHTML = `<a href="/r/${content}">Click</a>`;
            break;
        default:
            try {
                var blob = content as Blob;
                element.innerHTML = await blob.text();
            } catch (e) {
                element.innerHTML = content;
            }
    }
    element.hidden = false;
}
