import * as React from 'react';

import './Modal.css';


export interface ModalProps
{
    title:      string;
    onClose:    ( ) => any;
}


export class Modal extends React.Component<ModalProps>
{
    public render ( )
    {
        return (
            <div className='Modal-container'>
                <div className='Modal-backdrop' onClick={this.props.onClose}></div>

                <div className='Modal'>
                    <div className='Modal-header'>
                        <div className='Modal-title'>
                            {this.props.title}
                        </div>

                        <div className='Modal-close' onClick={this.props.onClose}>
                            X
                        </div>
                    </div>

                    <div className='Modal-content'>
                        {this.props.children}
                    </div>
                </div>
            </div>
        );
    }
}
