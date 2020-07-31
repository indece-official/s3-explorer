import * as React from 'react';

import './Button.css';


export interface ButtonProps
{
    type?:      'button' | 'reset' | 'submit';
    title?:     string;
    disabled?:  boolean;
    onClick?:   ( evt: any ) => any;
}


export class Button extends React.Component<ButtonProps>
{
    public render ( )
    {
        return (
            <button
                title={this.props.title}
                data-testid='Button'
                className={'Button' + (this.props.disabled ? ' disabled' : '')}
                type={this.props.type || 'button'}
                onClick={this.props.onClick}>
                {this.props.children}
            </button>
        );
    }
}
