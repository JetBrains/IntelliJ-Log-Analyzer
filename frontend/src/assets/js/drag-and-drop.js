//Drag and drop .zip to analyze
$(document).ready(function () {
    let zipWriter = new zip.ZipWriter(new zip.Data64URIWriter("application/zip"));
    let dropzone = $("body>div#dropzone")[0];
    let loader = $("#dropzone .loader");
    let disclamer = $("#dropzone .disclaimer");
    let directorySelector =  $("#select-dir");
    let archiveSelector = $("#select-archive")
    let fileUploader = $("#file-uploader");
    let fileAnalyzer = $("#file-analyzer");
    let lastTarget = null;

    //Adds all the content of a directory to zipWriter (object that contains zip archive to send to the backend)
    const scanDir = async (entry) => {
        return (new Promise((resolve) => {
            if (entry.isDirectory) {
                let directoryReader = entry.createReader();
                directoryReader.readEntries(function (entries) {
                    Promise.all(entries.map(scanDir)).then(resolve);
                });
            } else {
                entry.file(async file => {
                    await zipWriter.add(entry.fullPath, new zip.BlobReader(file));
                    resolve("success");
                })
            }

        }))
    }

    //check the file type and send it to the appropriate backend function
    const processFile = async (entry) => {
        let zipFile = entry.name.match(/\.zip/);
        let log = entry.name.match(/\.log.*/);
        return (new Promise((resolve) => {
            if (zipFile) {
                entry.file(async file => {
                    let reader = new FileReader();
                    reader.readAsDataURL(file);
                    reader.onload = await async function () {
                        resolve(await window.go.main.App.UploadArchive(reader.result));
                    }
                })
            } else if (log) {
                entry.file(async file => {
                    let reader = new FileReader();
                    reader.readAsDataURL(file);
                    reader.onload = await async function () {
                        resolve(await window.go.main.App.UploadLogFile(entry.name, reader.result));
                    }
                })
            } else {
                console.log("File type is not supported")
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
        loader.show();
        disclamer.hide();
        zipWriter = new zip.ZipWriter(new zip.Data64URIWriter("application/zip"));
        let items = e.dataTransfer.items

        let result = "";
        for (let i = 0; i < items.length; i++) {
            if (DataTransferItem.prototype.webkitGetAsEntry) {
                let entry = e.dataTransfer.items[i].webkitGetAsEntry();
                if (entry.isFile) {
                    result = await processFile(entry);
                } else if (entry.isDirectory) {
                    await scanDir(entry);
                    let zipFile = await zipWriter.close();
                    result = await window.go.main.App.UploadArchive(zipFile)
                }
            } else {
                console.log("webkitGetAsEntry is not supported")
            }
        }
        if (result) {
            fileUploader.hide();
            fileAnalyzer.show();
            render()
        }
        loader.hide();
        disclamer.show();
        dropzone.style.visibility = "hidden";
        dropzone.style.opacity = 0;
        return false;
    });
    window.addEventListener('dragover', function (ev) {
        ev.preventDefault();
    });
    directorySelector.on('click', async () => {
        var openedDir = await window.go.main.App.OpenFolder()
        if (openedDir.length > 0) {
            fileUploader.hide();
            fileAnalyzer.show();
            render()
        }
    })
    archiveSelector.on('click', async () => {
        var openedArchive = await window.go.main.App.OpenArchive()
        if (openedArchive.length > 0) {
            fileUploader.hide();
            fileAnalyzer.show();
            render()
        }
    })
})
