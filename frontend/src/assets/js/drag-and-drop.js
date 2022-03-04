//Drag and drop .zip to analyze
$(document).ready(function () {
    let zip = new JSZip();
    var dropzone = $("body>div#dropzone")[0];
    var lastTarget = null;
    async function scanDir(item) {
        return (new Promise((resolve, reject) => {
            if (item.isDirectory) {
                let directoryReader = item.createReader();
                directoryReader.readEntries(function (entries) {
                    Promise.all(entries.map(scanDir)).then(resolve);
                });
            } else {
                item.file(file => {
                    zip.file(item.fullPath, file);
                    resolve("success")
                })
            }

        }))
    }

    window.addEventListener('dragenter', function (ev) {
        lastTarget = ev.target;
        dropzone.style.visibility = ""
        dropzone.style.opacity = 1;

    });
    window.addEventListener('dragleave', function (ev) {
        if (ev.target === lastTarget || ev.target === document) {
            dropzone.style.visibility = "hidden";
            dropzone.style.opacity = 0;
        }
    });
    window.addEventListener('drop', async function (e) {
        e.preventDefault();
        e.stopPropagation();
        let items = e.dataTransfer.items

        const processFile = async (entry) => {
            let zipFile = entry.name.match(/\.zip/);
            let log = entry.name.match(/\.log.*/);
            if (zipFile) {
                entry.file(async file => {
                    debugger;
                    let content = await file.arrayBuffer()
                    result = await window.go.main.App.UploadArchive(content);
                    resolve(result)
                })
                //todo: pass file to backend

            } else if (log) {
                result = await window.go.main.App.UploadLogFile();
            } else {
                console.log("Fail ne podderjivaets")
            }
        }

        let result = "";
        for (let i = 0; i < items.length; i++) {
            let entry = e.dataTransfer.items[i].webkitGetAsEntry();
            if (entry.isFile) {
                await processFile(entry);
            } else if (entry.isDirectory) {
                await scanDir(entry)
                let zipFile = await zip.generateAsync({type: "array", compression: "DEFLATE",
                    compressionOptions: {
                        level: 1
                    }},function updateCallback(metadata) {
                    console.log("progression: " + metadata.percent.toFixed(2) + " %");
                    if(metadata.currentFile) {
                        console.log("current file = " + metadata.currentFile);
                    }
                })
                result = await window.go.main.App.UploadArchive(zipFile)
                console.log(result);
            }
        }
        if (result) {
            $("#file-uploader").hide();
            $("#file-analyzer").show();
            render()
        }

        dropzone.style.visibility = "hidden";
        dropzone.style.opacity = 0;
        return false;
    });

    window.addEventListener('dragover', function (ev) {
        ev.preventDefault();
    });
})
