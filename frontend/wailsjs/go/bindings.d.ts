import * as models from './models';

export interface go {
  "main": {
    "App": {
		CheckForUpdates():Promise<boolean>
		GetEntityNamesWithLineHighlightingColors():Promise<string>
		GetFilters():Promise<string>
		GetLogs():Promise<string>
		GetOtherFileContent(arg1:string):Promise<string>
		GetSetting(arg1:string):Promise<number>
		GetSettingsScreenHTML():Promise<string>
		GetStaticInfo():Promise<string>
		GetSummary():Promise<string>
		GetThreadDumpFileContent(arg1:string,arg2:string):Promise<string>
		GetThreadDumpsFilters(arg1:string):Promise<string>
		OpenArchive():Promise<string>
		OpenFolder():Promise<string>
		OpenIndexingReport(arg1:string):Promise<void>
		OpenIndexingSummaryForProject(arg1:string):Promise<void>
		RenderSystemMenu():Promise<void>
		SaveSetting(arg1:string,arg2:number):Promise<void>
		SetFilters(arg1:any):Promise<string>
		ShowNoUpdatesMessage():Promise<void>
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
