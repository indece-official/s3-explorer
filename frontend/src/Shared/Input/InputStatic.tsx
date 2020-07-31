import * as React from 'react';

import './Input.css';


export interface InputStaticProps
{
    label:      string;
}


export class InputStatic extends React.Component<InputStaticProps>
{
    public render ( )
    {
        return (
            <div className='Input'>
                <label>
                    <div
                        className='Input-label'>
                        {this.props.label}
                    </div>

                    <div className='Input-value'>
                        {this.props.children}
                    </div>
                </label>
            </div>
        );
    }
}
