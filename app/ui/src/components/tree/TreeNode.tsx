import {useState} from 'react'
import {useQuery} from '@tanstack/react-query'
import {
    ChevronDownIcon,
    ChevronRightIcon,
    CubeIcon,
    DocumentDuplicateIcon,
    FolderIcon,
} from '@heroicons/react/24/outline'
import {treeApi, TreeNode as TreeNodeType} from '@/services/treeService'
import {Loading} from '@/components/Loading'

interface TreeNodeProps {
    node: TreeNodeType
    level: number
    selectedId: string | null
    onSelect: (node: TreeNodeType) => void
    onContextMenu?: (node: TreeNodeType, e: React.MouseEvent) => void
}

export function TreeNode({node, level, selectedId, onSelect, onContextMenu}: TreeNodeProps) {
    const [isExpanded, setIsExpanded] = useState(false)
    const [offset, setOffset] = useState(0)
    const LIMIT = 25

    // Only fetch children when expanded
    const {data, isLoading, isFetching} = useQuery({
        queryKey: ['tree-children', node.id, offset],
        queryFn: () =>
            treeApi.getChildren({
                parent_id: node.id,
                limit: LIMIT,
                offset: offset,
            }),
        enabled: isExpanded,
        staleTime: 5 * 60 * 1000, // 5 minutes
    })

    // Track all loaded children across pagination
    const [allChildren, setAllChildren] = useState<TreeNodeType[]>([])

    // When new data arrives, append to all children
    if (data && data.items.length > 0) {
        const existingIds = new Set(allChildren.map((c) => c.id))
        const newChildren = data.items.filter((item) => !existingIds.has(item.id))
        if (newChildren.length > 0) {
            setAllChildren((prev) => [...prev, ...newChildren])
        }
    }

    const handleToggle = (e: React.MouseEvent) => {
        e.stopPropagation()
        if (!node.has_children) return
        setIsExpanded(!isExpanded)
    }

    const handleSelect = () => {
        onSelect(node)
    }

    const handleContextMenu = (e: React.MouseEvent) => {
        e.preventDefault()
        e.stopPropagation()
        onContextMenu?.(node, e)
    }

    const handleLoadMore = () => {
        setOffset((prev) => prev + LIMIT)
    }

    // Icon based on type
    const getIcon = () => {
        switch (node.o_type) {
            case 'folder':
                return <FolderIcon className="h-4 w-4 text-warning"/>
            case 'object':
                return <CubeIcon className="h-4 w-4 text-primary"/>
            case 'variant':
                return <DocumentDuplicateIcon className="h-4 w-4 text-secondary"/>
            default:
                return <CubeIcon className="h-4 w-4"/>
        }
    }

    const isSelected = selectedId === node.id
    const paddingLeft = level * 12 + 8 // 12px per level + 8px base

    return (
        <div>
            {/* Node itself */}
            <div
                className={`
                    flex items-center gap-1 px-2 py-1.5 cursor-pointer
                    hover:bg-base-200 transition-colors
                    ${isSelected ? 'bg-primary/10 border-l-2 border-primary' : ''}
                `}
                style={{paddingLeft: `${paddingLeft}px`}}
                onClick={handleSelect}
                onContextMenu={handleContextMenu}
            >
                {/* Expand/Collapse chevron */}
                <button
                    onClick={handleToggle}
                    className="btn btn-ghost btn-xs btn-square"
                    disabled={!node.has_children}
                >
                    {node.has_children ? (
                        isExpanded ? (
                            <ChevronDownIcon className="h-3 w-3"/>
                        ) : (
                            <ChevronRightIcon className="h-3 w-3"/>
                        )
                    ) : (
                        <span className="w-3"/> // Spacer for alignment
                    )}
                </button>

                {/* Icon */}
                {getIcon()}

                {/* Label */}
                <span className="text-sm truncate flex-1" title={node.o_key}>
                    {node.o_key}
                </span>

                {/* Children count badge */}
                {node.has_children && (
                    <span className="badge badge-xs badge-ghost">{node.children_count}</span>
                )}

                {/* Loading indicator */}
                {isExpanded && isLoading && <Loading variant="spinner" size="xs"/>}
            </div>

            {/* Children (recursive) */}
            {isExpanded && (
                <div>
                    {allChildren.map((child) => (
                        <TreeNode
                            key={child.id}
                            node={child}
                            level={level + 1}
                            selectedId={selectedId}
                            onSelect={onSelect}
                            onContextMenu={onContextMenu}
                        />
                    ))}

                    {/* Load More button */}
                    {data && data.has_more && (
                        <div style={{paddingLeft: `${(level + 1) * 12 + 8}px`}}>
                            <button
                                onClick={handleLoadMore}
                                className="btn btn-ghost btn-xs w-full justify-start gap-2 text-xs"
                                disabled={isFetching}
                            >
                                {isFetching ? (
                                    <>
                                        <Loading variant="spinner" size="xs"/>
                                        Loading...
                                    </>
                                ) : (
                                    <>
                                        Load {Math.min(LIMIT, data.total - allChildren.length)} more...
                                    </>
                                )}
                            </button>
                        </div>
                    )}
                </div>
            )}
        </div>
    )
}