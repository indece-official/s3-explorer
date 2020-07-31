import * as React from 'react';

import './Link.css';


export interface LinkProps
{
    url:    string;
    title?: string;
}


export class Link extends React.Component<LinkProps>
{
    constructor ( props: LinkProps )
    {
        super(props);

        this._onClick = this._onClick.bind(this);
    }


    private async _onClick ( ): Promise<void>
    {
        if ( typeof((window as any).systemOpenLink) === 'function' )
        {
            await (window as any).systemOpenLink(this.props.url);
        }
        else
        {
            window.open(this.props.url, '_blank');
        }
    }


    public render ( )
    {
        return (
            <span
                className='Link'
                onClick={this._onClick}>
                {this.props.title || this.props.url}
            </span>
        );
    }
}
