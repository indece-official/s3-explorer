import * as React from 'react';

import './Input.css';


export interface InputSelectOption
{
    label: string;
    value: string;
}


export interface InputSelectProps
{
    label:      string;
    name:       string;
    value:      string;
    options:    Array<InputSelectOption>;
    error?:     string | boolean;
    onChange?:  ( evt: any ) => any;
}


export class InputSelect extends React.Component<InputSelectProps>
{
    public render ( )
    {
        return (
            <div className='Input'>
                <label>
                    <div className='Input-value'>
                        <div className={'Input-input' + (this.props.error ? ' invalid' : '')}>
                            <select
                                value={this.props.value}
                                name={this.props.name}
                                onChange={this.props.onChange}>
                                {this.props.label ? 
                                    <option value=''>{this.props.label}</option>
                                : null}

                                {this.props.options.map( ( option ) => (
                                    <option
                                        key={option.value}
                                        value={option.value}>
                                        {option.label}
                                    </option>
                                ))}
                            </select>
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
