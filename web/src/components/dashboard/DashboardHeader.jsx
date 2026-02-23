/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React from 'react';
import { Button } from '@douyinfe/semi-ui';
import { RefreshCw, Search } from 'lucide-react';

const DashboardHeader = ({
  getGreeting,
  greetingVisible,
  showSearchModal,
  refresh,
  loading,
  t,
}) => {
  const ICON_BUTTON_CLASS = 'text-white hover:bg-opacity-80 !rounded-xl';

  return (
    <div className='flex items-center justify-between mb-4'>
      <h2
        className='text-2xl font-semibold claude-title transition-opacity duration-1000 ease-in-out'
        style={{ opacity: greetingVisible ? 1 : 0 }}
      >
        {getGreeting}
      </h2>
      <div className='flex gap-3'>
        <Button
          type='tertiary'
          icon={<Search size={16} />}
          onClick={showSearchModal}
          style={{
            background: 'linear-gradient(135deg, #D9A775 0%, #CC7A5F 100%)',
            border: 'none',
            color: 'white',
            borderRadius: '12px',
            width: '40px',
            height: '40px',
            padding: '0',
          }}
        />
        <Button
          type='tertiary'
          icon={<RefreshCw size={16} />}
          onClick={refresh}
          loading={loading}
          style={{
            background: 'linear-gradient(135deg, #E8D5C4 0%, #A67B5B 100%)',
            border: 'none',
            color: 'white',
            borderRadius: '12px',
            width: '40px',
            height: '40px',
            padding: '0',
          }}
        />
      </div>
    </div>
  );
};

export default DashboardHeader;
