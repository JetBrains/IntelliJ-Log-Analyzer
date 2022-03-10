export interface go {
  "main": {
    "App": {
		FilterUpdate(arg1:number):Promise<string>
		GetFilters():Promise<string>
		GetLogs():Promise<string>
		GetStaticInfo():Promise<string>
		OpenArchive():Promise<string>
		OpenFolder():Promise<string>
		SetFilters(arg1:any):Promise<string>
		UploadArchive(arg1:string):Promise<string>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
