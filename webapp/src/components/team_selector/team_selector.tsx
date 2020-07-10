import React, {useEffect, useState} from 'react';
import Select, {ActionTypes, ControlProps} from 'react-select';
import {Team} from 'mattermost-redux/types/teams';


interface Props {
    teamId?: string;
    getTeams: () => Promise<Team[]>;
    onSelectedChange: (teamId?: string) => void;
}

interface Option {
    value: string;
    label: JSX.Element;
    teamId: string;
}

interface ActionObj {
    action: ActionTypes;
}

export default function TeamSelector(props: Props) {
    const [teamOptions, setTeamOptions] = useState<Option[]>([]);

    async function fetchTeams() {
        const teams = await props.getTeams();
        const optionList = teams.map((team: Team) => {
            return ({
                value: team.display_name,
                label: (
                    <div>{team.display_name}</div>
                ),
                teamId: team.id,
            });
        });
        setTeamOptions(optionList);
    }

    // Fill in the userOptions on mount.
    useEffect(() => {
        fetchTeams();
    }, []);

    const [selected, setSelected] = useState<Option | null>(null);

    // Whenever the team changes we have to set the selected, but we can only do this once we
    // have teamOptions
    useEffect(() => {
        if (teamOptions === []) {
            return;
        }

        const team = teamOptions.find((option: Option) => option.teamId === props.teamId);
        if (team) {
            setSelected(team);
        } else {
            setSelected(null);
        }
    }, [teamOptions, props.teamId]);

    const onSelectedChange = async (value: Option | undefined, action: ActionObj) => {
        if (action.action === 'clear') {
            return;
        }
        if (value && selected && value.teamId === selected.teamId) {
            return;
        }
        if (value){
            props.onSelectedChange(value.teamId);
        } else {
            props.onSelectedChange(value);
        }
    };
    
    return (
        <Select
            isClearable={true}
            options={teamOptions}
            styles={selectStyles}
            defaultValue={selected}
            onChange={(option, action) => onSelectedChange(option as Option, action as ActionObj)}
            classNamePrefix='insights-team-select'
            className='insights-team-select'
        />
    );
}


// styles for the select component
const selectStyles = {
    control: (provided: any) => ({...provided, minWidth: 120, margin: 8}),
    menu: () => ({boxShadow: 'none'}),
    option: (provided: any, state: any) => {
        const hoverColor = 'rgba(20, 93, 191, 0.08)';
        const bgHover = state.isFocused ? hoverColor : 'transparent';
        return {
            ...provided,
            backgroundColor: state.isSelected ? hoverColor : bgHover,
            color: 'unset',
        };
    },
};