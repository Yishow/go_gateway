import React from 'react';
import { useNavigate } from 'react-router-dom';
import { MappingWizard, type WizardFormData } from '../../components/datalink/wizard';

/**
 * Mapping Wizard 頁面
 */
export const MappingWizardPage: React.FC = () => {
  const navigate = useNavigate();

  const handleComplete = (_data: WizardFormData) => {
    // 導航回 Mappings 頁面
    navigate('/datalink/mappings');
  };

  const handleCancel = () => {
    navigate(-1);
  };

  return (
    <MappingWizard
      onComplete={handleComplete}
      onCancel={handleCancel}
    />
  );
};

export default MappingWizardPage;
