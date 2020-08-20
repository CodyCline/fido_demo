import * as React from 'react';
import { ReactComponent as CaretUp } from '../../assets/caret-up-solid.svg';
import { ReactComponent as CaretDown } from '../../assets/caret-down-solid.svg';
import './spoiler.css';

export const Spoiler = ({ children, title } :any) => {
    const [toggled, toggleSpoiler] = React.useState<boolean>(false);
    const onToggle = () => {
        toggleSpoiler(!toggled);
    }
    return (
        <div className="spoiler">
            <ul onClick={onToggle} role="navigation" className="spoiler-header">
                <li>{title}</li>
                <li>
                    {toggled ? <CaretUp className="spoiler-icon"/>: <CaretDown className="spoiler-icon"/>}
                </li>
            </ul>
            {toggled &&
                <div className="spoiler-body">
                    {children}
                </div>
            }
        </div>
    );
}