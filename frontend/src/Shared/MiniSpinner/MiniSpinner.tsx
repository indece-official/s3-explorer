import * as React from 'react';

import './MiniSpinner.css';


export class MiniSpinner extends React.Component
{
    public render ( )
    {
        return (
            <div className='MiniSpinner'>
                <span className='MiniSpinner-dot1'>.</span>
                <span className='MiniSpinner-dot2'>.</span>
                <span className='MiniSpinner-dot3'>.</span>
            </div>
        );
    }
}