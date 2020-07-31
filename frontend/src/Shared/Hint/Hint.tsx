import * as React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faInfoCircle } from '@fortawesome/free-solid-svg-icons';


export interface HintProps
{
    text: string;
}


export class Hint extends React.Component<HintProps>
{
    public render ( )
    {
        return (
            <span
                className='Hint'
                title={this.props.text}>
                <FontAwesomeIcon
                    icon={faInfoCircle}
                />
            </span>
        )
    }
}
