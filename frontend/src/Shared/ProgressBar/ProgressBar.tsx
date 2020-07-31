import * as React from 'react';

import './ProgressBar.css';


export interface ProgressBarProps
{
    value: number;
    total: number;
}


export class ProgressBar extends React.Component<ProgressBarProps>
{
    public render ( )
    {
        const progress = Math.ceil(this.props.value / (this.props.total || 1) * 10000) / 100;

        return (
            <div className='ProgressBar'>
                <div
                    className='ProgressBar-progress'
                    style={{width: `${progress}%`}}
                />

                <div className='ProgressBar-label'>
                    {progress}%
                </div>
            </div>
        );
    }
}
