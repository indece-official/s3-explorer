import * as React from 'react';

import './Input.css';


export interface InputTextProps
{
    label:      string;
    name:       string;
    value:      string;
    error?:     string | boolean;
    onChange?:  ( evt: any ) => any;
}


export class InputText extends React.Component<InputTextProps>
{
    public render ( )
    {
        return (
            <div className='Input'>
                <label>
                    <div className='Input-value'>
                        <div className={'Input-input' + (this.props.error ? ' invalid' : '')}>
                            <input
                                type='text'
                                placeholder={this.props.label}
                                value={this.props.value}
                                name={this.props.name}
                                onChange={this.props.onChange}
                            />
                        </div>

                        {this.props.error && this.props.error !== true ?
                            <div className='Input-error'>
                                {this.props.error}
                            </div>
                        : null}

                        {this.props.children}
                    </div>
                </label>
            </div>
        );
    }
}
