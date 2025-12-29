import {useState} from 'react'
import {useQuery} from '@tanstack/react-query'
import {treeApi, TreeNode as TreeNodeType} from '@/services/treeService'
import {TreeNode} from './TreeNode'
import {Loading} from '@/components/Loading'
import {FolderIcon} from '@heroicons/react/24/outline'
import {AddEntityModal} from '@/components/tree/AddEntityModal'

// Root entity ID from your migrations
const ROOT_ENTITY_ID = '00000000-0000-0000-0000-000000000001'

interface ContextMenuState {
    node: TreeNodeType | 'root'
    x: number
    y: number
}

export function TreeView() {
    const [selectedNode, setSelectedNode] = useState<TreeNodeType | null>(null)
    const [contextMenu, setContextMenu] = useState<ContextMenuState | null>(null)
    const [isRootSelected, setIsRootSelected] = useState(false)
    const [isAddModalOpen, setIsAddModalOpen] = useState(false)
    const [addEntityParent, setAddEntityParent] = useState<{
        id: string | null
        path: string
    } | null>(null)

    // Fetch root's children to start
    const {data, isLoading, error} = useQuery({
        queryKey: ['tree-children', ROOT_ENTITY_ID, 0],
        queryFn: () =>
            treeApi.getChildren({
                parent_id: ROOT_ENTITY_ID,
                limit: 25,
                offset: 0,
            }),
        staleTime: 5 * 60 * 1000,
    })

    const handleNodeSelect = (node: TreeNodeType) => {
        setSelectedNode(node)
        setIsRootSelected(false)
        setContextMenu(null)
        console.log('Selected node:', node)
    }

    const handleRootSelect = () => {
        setSelectedNode(null)
        setIsRootSelected(true)
        setContextMenu(null)
        console.log('Selected root folder')
    }

    const handleContextMenu = (node: TreeNodeType | 'root', e: React.MouseEvent) => {
        e.preventDefault()
        e.stopPropagation()
        setContextMenu({
            node,
            x: e.clientX,
            y: e.clientY,
        })
    }

    const handleCloseContextMenu = () => {
        setContextMenu(null)
    }

    const handleAddEntity = () => {
        if (!contextMenu) return

        // Determine parent info based on context menu node
        if (contextMenu.node === 'root') {
            setAddEntityParent({
                id: ROOT_ENTITY_ID,
                path: '/',
            })
        } else {
            setAddEntityParent({
                id: contextMenu.node.id,
                path: contextMenu.node.o_path,
            })
        }

        setIsAddModalOpen(true)
        setContextMenu(null)
    }

    const handleRemoveEntity = () => {
        // TODO: Implement remove entity logic
        if (contextMenu?.node !== 'root') {
            console.log('Remove entity:', contextMenu?.node)
        }
        setContextMenu(null)
    }

    const handleAddModalClose = () => {
        setIsAddModalOpen(false)
        setAddEntityParent(null)
    }

    // Close context menu on outside click
    if (contextMenu) {
        const handleClickOutside = () => handleCloseContextMenu()
        document.addEventListener('click', handleClickOutside, {once: true})
    }

    if (isLoading) {
        return (
            <div className="flex h-full items-center justify-center">
                <Loading variant="spinner" size="md"/>
            </div>
        )
    }

    if (error) {
        return (
            <div className="p-4 text-center text-error text-sm">
                Failed to load tree
            </div>
        )
    }

    const isRootNode = contextMenu?.node === 'root'

    return (
        <div className="relative h-full overflow-y-auto">
            <div className="py-2">
                {/* Root Folder Node */}
                <div
                    className={`
                        flex items-center gap-1 px-2 py-1.5 cursor-pointer
                        hover:bg-base-200 transition-colors
                        ${isRootSelected ? 'bg-primary/10 border-l-2 border-primary' : ''}
                    `}
                    style={{paddingLeft: '8px'}}
                    onClick={handleRootSelect}
                    onContextMenu={(e) => handleContextMenu('root', e)}
                >
                    {/* Empty spacer for alignment with children */}
                    <span className="w-6"/>

                    {/* Folder Icon */}
                    <FolderIcon className="h-4 w-4 text-warning"/>

                    {/* Label */}
                    <span className="text-sm truncate flex-1 font-semibold">
                        System Root
                    </span>

                    {/* Children count badge */}
                    {data && data.total > 0 && (
                        <span className="badge badge-xs badge-ghost">{data.total}</span>
                    )}
                </div>

                {/* Root's children */}
                {data && data.items.length > 0 && (
                    <div>
                        {data.items.map((node) => (
                            <TreeNode
                                key={node.id}
                                node={node}
                                level={1}
                                selectedId={selectedNode?.id || null}
                                onSelect={handleNodeSelect}
                                onContextMenu={(node, e) => handleContextMenu(node, e)}
                            />
                        ))}
                    </div>
                )}

                {/* Empty state under root */}
                {data && data.items.length === 0 && (
                    <div className="pl-8 py-4 text-center text-base-content/50 text-xs">
                        No entities yet. Right-click System Root to add one.
                    </div>
                )}
            </div>

            {/* Context Menu */}
            {contextMenu && (
                <div
                    className="fixed z-50 menu bg-base-100 rounded-box shadow-lg w-48 p-2"
                    style={{
                        left: `${contextMenu.x}px`,
                        top: `${contextMenu.y}px`,
                    }}
                >
                    <li>
                        <button onClick={handleAddEntity}>
                            <span className="text-xs">Add Entity</span>
                        </button>
                    </li>
                    {/* Only show remove option for non-root nodes */}
                    {!isRootNode && (
                        <li>
                            <button onClick={handleRemoveEntity} disabled className="text-error">
                                <span className="text-xs opacity-50">Remove (TODO)</span>
                            </button>
                        </li>
                    )}
                </div>
            )}

            {/* Add Entity Modal */}
            {isAddModalOpen && addEntityParent && (
                <AddEntityModal
                    parentId={addEntityParent.id}
                    parentPath={addEntityParent.path}
                    onClose={handleAddModalClose}
                />
            )}
        </div>
    )
}