import * as React from 'react';

import './NavBar.css';


export interface NavBarProps
{
    onSelectProfile:    ( ) => any;
    onAbout:            ( ) => any;
    onAddObject:        ( ) => any;
}


export class NavBar extends React.Component<NavBarProps>
{
    public render ( )
    {
        return (
            <div className='NavBar'>
                <div
                    className='NavBar-item'
                    onClick={this.props.onSelectProfile}>
                    Profiles
                </div>
                
                <div
                    className='NavBar-item'
                    onClick={this.props.onAbout}>
                    About
                </div>

                <div className='NavBar-spacer'></div>

                <div
                    className='NavBar-item'
                    onClick={this.props.onAddObject}>
                    Upload a file
                </div>
            </div>
        );
    }
}
