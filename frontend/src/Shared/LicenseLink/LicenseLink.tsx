import * as React from 'react';

import './LicenseLink.css';


export class LicenseLink extends React.Component
{
    constructor ( props: any )
    {
        super(props);

        this._onClick = this._onClick.bind(this);
    }


    private async _onClick ( ): Promise<void>
    {
        if ( typeof((window as any).showLicense) === 'function' )
        {
            await (window as any).showLicense();
        }
        else
        {
            window.open('/LICENSE.txt', '_blank');
        }
    }


    public render ( )
    {
        return (
            <span
                className='LicenseLink'
                onClick={this._onClick}>
                View the license
            </span>
        );
    }
}
