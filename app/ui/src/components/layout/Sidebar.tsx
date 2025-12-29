import {TreeView} from "@/components/tree/TreeView.tsx";

export function Sidebar() {
    return (
        <aside className="min-h-full w-64 bg-base-300/20 text-base-content border-r border-base-300">
            {/* Header */}
            <div className="px-4 py-3 border-b border-base-300">
                <h2 className="text-sm font-semibold text-base-content/70">Entities</h2>
            </div>

            {/* Tree View */}
            <div className="h-[calc(100vh-8rem)] overflow-hidden">
                <TreeView/>
            </div>
        </aside>
    )
}