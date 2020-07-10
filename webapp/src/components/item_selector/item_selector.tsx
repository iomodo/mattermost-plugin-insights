import React, {useEffect, useState} from 'react';
import Select, {ActionTypes} from 'react-select';
import {Item} from '../../types/posts';


interface Props {
    id?: string;
    getItems: () => Promise<Item[]>;
    onSelectedChange: (id?: string) => void;
}

interface Option {
    value: string;
    label: JSX.Element;
    id: string;
}

interface ActionObj {
    action: ActionTypes;
}

export default function ItemSelector(props: Props) {
    const [itemOptions, setItemOptions] = useState<Option[]>([]);

    async function fetchItems() {
        const items = await props.getItems();
        const optionList = items.map((item: Item) => {
            return ({
                value: item.display_name,
                label: (
                    <div>{item.display_name}</div>
                ),
                id: item.id,
            });
        });
        setItemOptions(optionList);
    }

    // Fill in the userOptions on mount.
    useEffect(() => {
        fetchItems();
    }, []);

    const [selected, setSelected] = useState<Option | null>(null);

    // Whenever the item changes we have to set the selected, but we can only do this once we
    // have itemOptions
    useEffect(() => {
        if (itemOptions === []) {
            return;
        }

        const item = itemOptions.find((option: Option) => option.id === props.id);
        if (item) {
            setSelected(item);
        } else {
            setSelected(null);
        }
    }, [itemOptions, props.id]);

    const onSelectedChange = async (value: Option | undefined, action: ActionObj) => {
        if (action.action === 'clear') {
            return;
        }
        if (value && selected && value.id === selected.id) {
            return;
        }
        if (value){
            props.onSelectedChange(value.id);
        } else {
            props.onSelectedChange(value);
        }
    };
    
    return (
        <Select
            isClearable={true}
            options={itemOptions}
            styles={selectStyles}
            defaultValue={selected}
            onChange={(option, action) => onSelectedChange(option as Option, action as ActionObj)}
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