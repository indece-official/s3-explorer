import * as React from 'react';

import './Input.css';


export interface InputFileProps
{
    label:      string;
    name:       string;
    error?:     string | boolean;
    onChange?:  ( evt: any ) => any;
}


export class InputFile extends React.Component<InputFileProps>
{
    public render ( )
    {
        return (
            <div className='Input'>
                <label>
                    <div className='Input-value'>
                        <div className={'Input-input' + (this.props.error ? ' invalid' : '')}>
                            <input
                                type='file'
                                placeholder={this.props.label}
                                name={this.props.name}
                                onChange={this.props.onChange}
                                multiple
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
