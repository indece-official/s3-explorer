import { BackendService } from './BackendService';


export interface UpdateV1
{
    version:    string;
    datetime:   string;
}


export class S3UpdateService
{
    private static _instance: S3UpdateService;


    public static getInstance ( ): S3UpdateService
    {
        if ( ! this._instance )
        {
            this._instance = new S3UpdateService();
        }

        return this._instance;
    }


    private readonly _backendService: BackendService;


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
    }


    public async getUpdateAvailable ( ): Promise<UpdateV1 | null>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/update`,
            {
                method:     'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.update;
    }
}
