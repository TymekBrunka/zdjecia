wstaw = document.querySelectorAll("#dodaj, .dodaj")
ja = document.querySelector("ja > div:nth-child(2)")
ty = document.querySelector("ty > div:nth-child(2)")
pliki = {}
plikID = 0

pliki2 = {}
plikID2 = 0

debug = []

input = document.createElement("input")
input.type = "file"
input.setAttribute("multiple" ,"")

var ws = io();
ws.emit("wsc", {})

var jsz = new JSZip()

ws.on("setid", (msg) => {
    _ID = msg.id;
    document.querySelector("ja > div:first-child").innerText = `Ty (id: ${_ID})`
})

for (el of wstaw){
    el.onclick = () => {
        input.click()
    }
}

input.onchange = () => {
    files = []
    for (f of input.files){
        pliki[plikID] = f
        el = document.createElement("div")
        el.className = "karta"
        el.setAttribute("plikID", plikID)
        if (f.type.substring(0,5) === "image"){
            src = URL.createObjectURL(f)
            files.push({src: f, name: f.name, prev: false, type: f.type, id: _ID})
        } else {
            src = "static/img/plik.svg"
            // console.log(f.type.substring(0,5))
            files.push({src: src, name: f.name, prev: true, type: "image/svg", id: _ID})
        }
        el.innerHTML = `
            <img src="${src}" alt="${f.name}">
            ${f.name}
            <button class="del" onclick="usun(${plikID})">x</button>
        `
        ja.appendChild(el)
        plikID += 1
    }
    ws.emit("podglad", {f: files, id: _ID})
}

ws.on("podglad", (msg) => {
    console.log(msg.id)
    if (msg.id != _ID){
        console.log(msg)
        debug = msg
        for (f of msg.f){
            pliki2[plikID2] = f
            el = document.createElement("div")
            el.className = "karta"
            el.setAttribute("plikID", plikID)
            if (f.prev != true){
                src = receive_data(f.src, f.type)
            } else {
                src = "static/img/plik.svg"
            }
            el.innerHTML = `
                <img src="${src}" alt="${f.name}">
                ${f.name}
            `
            ty.appendChild(el)
            plikID2 += 1
        }
    }
})

ws.on("delete", (msg) => {
    document.querySelector(`[plikid="${msg.idp}"].karta`).remove()
    delete(pliki2[msg.id])
})

function receive_data(data, type) {
    var arrayBufferView = new Uint8Array(data);
    var blob = new Blob( [ arrayBufferView ], { type: type } );

    var img_url = URL.createObjectURL(blob);
    return img_url
}

function usun(id){
    document.querySelector(`[plikid="${id}"].karta`).remove()
    delete(pliki[id])
    ws.emit("delete", {idp: id, id: _ID})
}
