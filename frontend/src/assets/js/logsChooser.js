const dropzone = $("body>div#dropzone").first();
const loader = dropzone.find(".loader");
const disclamer = dropzone.find(".disclaimer");
const directorySelector = $("#select-dir");
const archiveSelector = $("#select-archive")
const fileUploader = $("#file-uploader");
const fileAnalyzer = $("#file-analyzer");
const IdeSelector = $("#select-running-ide")

$(document).ready(async function () {
    let zipWriter = new zip.ZipWriter(new zip.Data64URIWriter("application/zip"));
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

    //inserts the list of running and installed IDEs into #select-running-ide .dropdown
    IdeSelector.find(".options").first().html(await window.go.main.App.GetRunningIDEsDropdownHTML())
    window.addEventListener('dragenter', function (ev) {
        lastTarget = ev.target;
        dropzone.css('visibility', 'visible');
        dropzone.css('opacity', '1');

    });
    window.addEventListener('dragleave', function (ev) {
        if (ev.target === lastTarget || ev.target === document) {
            dropzone.css('visibility', 'hidden');
            dropzone.css('opacity', '0');
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
        dropzone.css('visibility', 'hidden');
        dropzone.css('opacity', '0');
        return false;
    });
    window.addEventListener('dragover', function (ev) {
        ev.preventDefault();
    });
    directorySelector.on('click', async () => {
        let path = await window.go.main.App.OpenFolder()
        initLogDirectory(path)
    })
    archiveSelector.on('click', async () => {
        let path = await window.go.main.App.OpenArchive()
        initLogDirectory(path)
    })
    IdeSelector.find(".button").first().on('click', async function () {
        let path = IdeSelector.find("li.active").attr("target");
        $(this).html("Loading...");
        await initLogDirectory(path)
        window.go.main.App.EnableLogsLiveUpdate()
    })
})
document.addEventListener('DOMContentLoaded', function () {
    window.runtime.EventsOn("LogsUpdated", function (s) {
        appendToMainEditor(s)
    })
})
async function initLogDirectory(path) {
    let openedLogsDir = await window.go.main.App.InitLogDirectory(path)
    if (openedLogsDir.length > 0) {
        fileUploader.hide();
        fileAnalyzer.show();
        render()
    }
}
