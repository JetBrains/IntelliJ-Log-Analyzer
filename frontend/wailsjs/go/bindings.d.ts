export interface go {
  "main": {
    "App": {
		FilterUpdate(arg1:number):Promise<string>
		GetEntityNamesWithLineHighlightingColors():Promise<string>
		GetFilters():Promise<string>
		GetLogs():Promise<string>
		GetStaticInfo():Promise<string>
		GetSummary():Promise<string>
		GetThreadDumpFileContent(arg1:string,arg2:string):Promise<string>
		GetThreadDumpsFilters(arg1:string):Promise<string>
		OpenArchive():Promise<string>
		OpenFolder():Promise<string>
		SetFilters(arg1:any):Promise<string>
		UploadArchive(arg1:string):Promise<string>
		UploadLogFile(arg1:string,arg2:string):Promise<string>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
