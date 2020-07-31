/* https://github.com/g6ling/React-Native-Tips/tree/master/How_to_upload_photo%2Cfile_in%20react-native */

export function futch ( url: string, opts: any = {}, onProgress?: ( evt: ProgressEvent ) => any ): Promise<any>
{
    return new Promise( ( res, rej ) =>
    {
        const xhr = new XMLHttpRequest();

        xhr.open(opts.method || 'get', url);

        for ( var k in opts.headers || {} )
        {
            xhr.setRequestHeader(k, opts.headers[k]);
        }

        xhr.onload = ( e ) => res(e.target);
        xhr.onerror = rej;

        if ( xhr.upload && typeof(onProgress) === 'function' )
        {
            xhr.upload.onprogress = onProgress; // event.loaded / event.total * 100 ; //event.lengthComputable
        }

        xhr.send(opts.body);
    });
}

