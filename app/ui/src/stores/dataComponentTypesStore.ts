export enum DataComponentCategory {
    Numeric = 'numeric',
    Text = 'text',
    Date = 'date',
}

export enum DataComponentIcon {
    Hash = 'hash',
    Calendar = 'calendar',
    Cursor = 'text-cursor'
}


export enum DataComponentType {
    Date = 'date',
    Varchar = 'varchar',
    Integer = 'integer',
}

export interface DataComponent {
    id: string
    label: string
    categories: DataComponentCategory
    tooltip: string
    icon: DataComponentIcon
    defaultDBType: DataComponentType
}
