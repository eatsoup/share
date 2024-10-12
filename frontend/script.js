document.addEventListener('paste', async (e) => {
    e.preventDefault();
    if (e.clipboardData.files?.length > 0) {
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


async function postIt(url, data) {
    const request = new XMLHttpRequest();
    request.onreadystatechange = function() {
        if (request.readyState == 4 && (request.status>=200 && request.status<400)) {
            setContent(`<a href="?id=${request.responseText}">click</a>`);
        }
    }
    request.open("POST", url);
    request.send(data);
}

async function getIt(url) {
    const request = new XMLHttpRequest();
    request.onreadystatechange = function() {
        if (request.readyState === 4 && (request.status>=200 && request.status<400)) {
            setContent(request.responseText, request.getResponseHeader("Content-Type"));
        }
    }
    request.open("GET", url);
    request.send();
}

async function showContent(id) {
    const data = await getIt(`/get?id=${id}`);
    if (data != undefined) setContent(data);
}

function setContent(content, ctype) {
    console.log("setting content")
    const element = document.getElementById("content");
    switch(ctype) {
        case "text/html":
            element.innerHTML = content;
            break;
        case "image/png":
            const blob = new Blob([content], {type: ctype});
            const reader = new FileReader();
            reader.onloadend = () => {
                const bstring = reader.result;
                element.innerHTML = `<img src="${bstring}">`;
            }
            reader.readAsDataURL(blob);
            break;
        default:
            element.innerHTML = content;
    }
    element.hidden = false;
}
