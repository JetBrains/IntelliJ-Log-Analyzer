export interface go {
  "main": {
    "App": {
		GetEntityNamesWithLineHighlightingColors():Promise<string>
		GetFilters():Promise<string>
		GetLogs():Promise<string>
		GetOtherFileContent(arg1:string):Promise<string>
		GetStaticInfo():Promise<string>
		GetSummary():Promise<string>
		GetThreadDumpFileContent(arg1:string,arg2:string):Promise<string>
		GetThreadDumpsFilters(arg1:string):Promise<string>
		OpenArchive():Promise<string>
		OpenFolder():Promise<string>
		OpenIndexingReport(arg1:string):Promise<void>
		OpenIndexingSummaryForProject(arg1:string):Promise<void>
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
