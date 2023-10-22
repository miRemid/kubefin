import { Route, Routes } from "react-router-dom";
import { AllClustersDashboard } from "./pages/dashboard/components/all-clusters/AllClustersDashboard";
import { SingleClusterDashboard } from "./pages/dashboard/components/single-cluster/SingleClusterDashboard";
import "@cloudscape-design/global-styles/index.css"
import { ClusterCostDashboard } from "./pages/dashboard/components/cluster-cost/ClusterCostDashboard";
import { WorkloadCostDashboard } from "./pages/dashboard/components/workload-cost/WorkloadCostDashboard";
import { NamespaceCostDashboard } from "./pages/dashboard/components/namespace-cost/NameSpaceCostDashboard";

export function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<AllClustersDashboard />} />
        <Route path="/dashboard">
          <Route index element={<AllClustersDashboard />} />
          <Route path=":clusterId" element={<SingleClusterDashboard />} />
        </Route>
        <Route path="/cost">
          <Route path=":clusterId/:clusterName/cluster" element={<ClusterCostDashboard />} />
          <Route path=":clusterId/:clusterName/workload" element={<WorkloadCostDashboard />} />
          <Route path=":clusterId/:clusterName/namespace" element={<NamespaceCostDashboard />} />
        </Route>
        {/* <Route path="*" element={<NotFound />} /> */}
      </Routes>
    </>
  );
}
export default App;
