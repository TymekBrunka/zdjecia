fetch('/ipv4', {
    method: 'GET',
    headers: {
        'Content-Type': 'text/plain'
    },
    responseType: 'text' // Explicitly set the responseType to 'text'
}).then((res) => {

    if (!res.ok) {
        console.log("server skill issue for ipv4 response");
    }
    return res.text();

}).then((data) => {

    document.querySelector("adres").innerHTML = `http://${data}:8000`;

}).catch(error => {
    console.error('There was a problem with ipv4 fetch operation:', error);
});

//-----functions

function toIMG(data, type) {
    let arrayBufferView
    data.arrayBuffer().then((res) => { arrayBufferView = res; console.log(res) });
    var blob = new Blob([arrayBufferView], { type: type });

    var img_url = URL.createObjectURL(blob);
    return img_url
}

class sender {
    static readFileAsDataURL(file) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = (event) => {
                const fileData = event.target.result;
                resolve(fileData);
            };

            reader.onerror = (error) => {
                reject(error);
            };

            reader.readAsDataURL(file);
        });
    }

    static genFIlesFromB64(files) {
        let newFiles = files.map((fob) => {
            let datastr = fob.data
            datastr = datastr.split(",")[1]
            let bytestr = atob(datastr)
            let bytenums = new Array(bytestr.length)
            for (let i = 0; i < bytestr.length; i++) {
                bytenums[i] = bytestr.charCodeAt(i)
            }
            let bytes = new Uint8Array(bytenums)
            let blob = new Blob([bytes], { type: fob.ft })
            let file = new File([blob], fob.name, { type: fob.ft })

            return { id: fob.id, file: file }
        })
        return newFiles
    }

    static downloadFile(file) {
        const blob = new Blob([file], { type: file.type });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = file.name;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
    }

    static async sendImg(files) {
        const newFiles = files.map(async (e) => {
            const fileData = await this.readFileAsDataURL(e.file);
            return {
                id: e.id,
                ft: e.file.type,
                size: e.file.size,
                name: e.file.name,
                data: fileData
            };
        });

        const resolvedFiles = await Promise.all(newFiles);

        fetch("/sendfiles",

            {
                method: 'POST',
                // responseType: 'text', // Explicitly set the responseType to 'text'
                body:
                    JSON.stringify({
                        // type: 0,
                        nr: nr,
                        files: resolvedFiles
                    }),
                headers: {
                    'Content-Type': 'application/json'
                },
            },
        );

        console.log("sending: ")
        console.log({
            // type: 0,
            nr: nr,
            files: resolvedFiles
        })
    }
}


function usun(id) {
    // let rm_element = 0;
    for (el of ja_grid.querySelectorAll("button[iid]").values()) {
        if (el.getAttribute("iid") == id) {
            el.remove();
            break;
        }
    }
    sock.send(JSON.stringify({
        type: 1,
        nr: nr,
        ids: [{
            id: id,
            name: files.find((e) => { return e.id == id }).file.name,
        }]
    }))
    files.splice(id, 1)
}

//-----main
nr = 0
sock = new WebSocket(`ws://${location.host}/ws`)

sock.onmessage = (e) => {
    console.log(e)
    msg = JSON.parse(e.data)
    switch (msg.type) {
        case (69420): {
            nr = msg.id

            newFiles = sender.genFIlesFromB64(msg.data[nr == 0 ? 0 : 1].files)
            files = files.concat(newFiles)
            dodaj_karty(newFiles, ja_grid, false)

            newFiles = sender.genFIlesFromB64(msg.data[nr == 0 ? 1 : 0].files)
            opfiles = opfiles.concat(newFiles)
            dodaj_karty(newFiles, ty_grid, false)

            break;
        }
        case (2137): {
            fetch("/givfilesplz",

                {
                    method: 'POST',
                    responseType: 'text', // Explicitly set the responseType to 'text'
                    body:
                        // JSON.stringify({
                        // type: 0,
                        nr,
                    // files: resolvedFiles
                    // }),
                    headers: {
                        'Content-Type': 'text/plain'
                    },
                },
            ).then((res) => {

                if (!res.ok) {
                    console.log("Server is mean, he refuses to give files my friend sent me");
                }
                return res.text();

            }).then((data) => {
                // console.log(data)
                data = JSON.parse(data)
                // console.log(data)
                newFiles = sender.genFIlesFromB64(data.files)
                opfiles = opfiles.concat(newFiles)
                dodaj_karty(newFiles, ty_grid, false)
            }).catch(error => {
                console.error('fetch was wrong: ', error);
            });
            break;
        }
        case (1): {
            for (karta of msg.ids) {
                opfiles.splice(karta.id, 1)
                for (el of ty_grid.querySelectorAll("button[iid]").values()) {
                    if (el.getAttribute("iid") == karta.id) {
                        el.remove();
                        break;
                    }
                }
            }
            break;
        }
    }
}

sock.onopen = () => {
    sock.send(`{"type" : 69, "id": 0, "ids": []}`)
}

_id = 0;
files = []
opfiles = []
fs_menu = document.createElement("input")
fs_menu.type = "file"
fs_menu.multiple = true
ja_grid = document.querySelector("ja > div:nth-child(2)")
ty_grid = document.querySelector("ty > div:nth-child(2)")

document.querySelector(".dodaj").onclick = () => {
    fs_menu.click()
}

function dodaj_karty(newfiles, gridbox, usuwalne) {

    for (file of newFiles) {
        karta = document.createElement("button")
        karta.className = "karta"
        karta.setAttribute("iid", file.id)
        karta.innerHTML = `
            <img src="${file.file.type.substring(0, 5) === "image" ? URL.createObjectURL(file.file) : "./img/plik.svg"}" alt="${file.file.name}">
            ${file.file.name}
            ${usuwalne ? '<button class="del">x</button>' : ''}
        `
        // const eid = _id //a to tak w razie czego jeśli _id jako argument dla funkcji by miał się zmienić
        // karta.onclick = () => {usun(eid)}
        if (usuwalne) { karta.setAttribute("onclick", `usun(${file.id})`); }
        gridbox.appendChild(karta)
    }
}

fs_menu.onchange = async (e) => {
    // console.log(e)
    i = 0;
    newFiles = Array.prototype.slice.call(fs_menu.files).map((e) => { return { id: -1, file: e }; })
    files = files.concat(newFiles)
    files = files.map((el) => { if (el.id == -1) { el.id = _id++; } return el; })
    console.log(newFiles)
    sender.sendImg(newFiles)
    dodaj_karty(newFiles, ja_grid, true)
    // sender.sendImg(newFiles)
}

document.getElementById("sciagnij").onclick = () => {
    for (file of opfiles) {
        sender.downloadFile(file.file)
    }
}

