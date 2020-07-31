import * as React from 'react';

import './Input.css';


export interface InputCheckboxProps
{
    label:      string;
    name:       string;
    value:      boolean;
    error?:     boolean;
    disabled?:  boolean;
    onChange?:  ( evt: any ) => any;
}


export class InputCheckbox extends React.Component<InputCheckboxProps>
{
    public render ( )
    {
        return (
            <div className='Input Input-checkbox'>
                <label className={'Input-checkbox-label' + (this.props.error ? ' invalid' : '')}>
                    <input
                        type='checkbox'
                        name={this.props.name}
                        checked={this.props.value}
                        disabled={this.props.disabled}
                        onChange={this.props.onChange}
                    />

                    <div className='Input-label'>
                        {this.props.label}
                    </div>

                    {this.props.children}
                </label>
            </div>
        );
    }
}
