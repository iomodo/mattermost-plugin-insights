export interface Item {
    id: string;
    display_name: string;
}

export interface Team extends Item {
    name: string;
}

export interface Channel extends Item {
    name: string;
}