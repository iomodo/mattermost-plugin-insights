export interface Item {
    id: string;
    display_name: string;
}

export interface Team extends Item {
    name: string;
}
