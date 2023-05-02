import { futch } from './futch';


export enum BackendMethod
{
    GET     = 'GET',
    POST    = 'POST'
}


export class BackendService
{
    private static _instance: BackendService;


    public static getInstance ( ): BackendService
    {
        if ( ! this._instance )
        {
            this._instance = new BackendService();
        }

        return this._instance;
    }


    private async getSessionToken ( ): Promise<string>
    {
        if ( typeof((window as any).s3SessionToken) !== 'function' )
        {
            return '';
        }

        return await (window as any).s3SessionToken();
    }


    public async fetchJson ( input: RequestInfo, init?: RequestInit ): Promise<any>
    {
        init = init || {};

        init.headers = {
            Authorization: `Bearer ${await this.getSessionToken()}`,
            ...(init.headers || {})
        };

        const resp = await fetch(input, {
            credentials: 'include',
            ...init
        });

        if ( ! resp.ok )
        {
            let text = `Can't send request to server: ${resp.status} ${resp.statusText}`;
            
            try
            {
                text = (await resp.text()) || text;
            }
            catch ( err ) { }

            throw new Error(text);
        }

        return await resp.json();
    }


    public async uploadFile ( path: string,
                              init: RequestInit,
                              onProgress?: ( evt: any ) => any ): Promise<any>
    {
        init = init || {};

        init.headers = {
            Authorization: `Bearer ${await this.getSessionToken()}`,
            ...(init.headers || {})
        };

        const resp = await futch(path, {
            ...init
        }, onProgress);

        return JSON.parse(resp.responseText);
    }
}
