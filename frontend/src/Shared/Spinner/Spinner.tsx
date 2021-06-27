import * as React from 'react';

import './Spinner.css';


export interface SpinnerProps
{
    active: boolean;
}


export class Spinner extends React.Component<SpinnerProps>
{
    public render ( )
    {
        if ( ! this.props.active )
        {
            return null;
        }

        return (
            <div className='Spinner'>
            </div>
        );
    }
}
