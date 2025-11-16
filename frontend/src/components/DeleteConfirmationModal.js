import React, { useState } from 'react';
import { Modal, Button, Input } from './ui';
import { AlertTriangle } from 'lucide-react';

function DeleteConfirmationModal({ 
  isOpen, 
  onClose, 
  onConfirm, 
  title = "Confirm Deletion",
  entityName,
  entityType,
  warningMessage 
}) {
  const [reason, setReason] = useState('');
  const [attestation, setAttestation] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!reason.trim()) {
      alert('Please provide a reason for deletion');
      return;
    }

    if (!attestation.trim()) {
      alert('Please provide your attestation');
      return;
    }

    setLoading(true);
    try {
      await onConfirm({ reason, attestation });
      handleClose();
    } catch (error) {
      console.error('Delete failed:', error);
      setLoading(false);
    }
  };

  const handleClose = () => {
    setReason('');
    setAttestation('');
    setLoading(false);
    onClose();
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      title={title}
      size="md"
    >
      <form onSubmit={handleSubmit}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          {/* Warning Banner */}
          <div style={{
            padding: '16px',
            background: 'var(--danger-light)',
            border: '1px solid var(--danger)',
            borderRadius: '8px',
            display: 'flex',
            gap: '12px',
            alignItems: 'flex-start'
          }}>
            <AlertTriangle size={24} color="var(--danger)" style={{ flexShrink: 0, marginTop: '2px' }} />
            <div>
              <h4 style={{ margin: '0 0 8px 0', color: 'var(--danger)', fontSize: '15px', fontWeight: '600' }}>
                Warning: This action will soft delete the record
              </h4>
              <p style={{ margin: 0, fontSize: '14px', color: 'var(--text-primary)' }}>
                {warningMessage || `You are about to delete ${entityName}. This action can be reversed by an administrator, but it will be recorded in the audit log.`}
              </p>
            </div>
          </div>

          {/* Entity Info */}
          {entityName && (
            <div style={{
              padding: '12px',
              background: 'var(--background-hover)',
              borderRadius: '8px',
              fontSize: '14px'
            }}>
              <div style={{ display: 'grid', gridTemplateColumns: '100px 1fr', gap: '8px' }}>
                <span style={{ color: 'var(--text-secondary)', fontWeight: '500' }}>Type:</span>
                <span style={{ fontWeight: '600' }}>{entityType}</span>
                <span style={{ color: 'var(--text-secondary)', fontWeight: '500' }}>Name:</span>
                <span style={{ fontWeight: '600' }}>{entityName}</span>
              </div>
            </div>
          )}

          {/* Reason Input */}
          <div>
            <label style={{
              display: 'block',
              marginBottom: '8px',
              fontSize: '14px',
              fontWeight: '500',
              color: 'var(--text-primary)'
            }}>
              Reason for Deletion <span style={{ color: 'var(--danger)' }}>*</span>
            </label>
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Enter a detailed reason for deleting this record..."
              rows={3}
              required
              style={{
                width: '100%',
                padding: '12px',
                border: '1px solid var(--border)',
                borderRadius: '8px',
                fontSize: '14px',
                fontFamily: 'inherit',
                resize: 'vertical'
              }}
            />
            <p style={{ margin: '6px 0 0 0', fontSize: '12px', color: 'var(--text-secondary)' }}>
              This reason will be permanently recorded in the audit log
            </p>
          </div>

          {/* Attestation Input */}
          <div>
            <label style={{
              display: 'block',
              marginBottom: '8px',
              fontSize: '14px',
              fontWeight: '500',
              color: 'var(--text-primary)'
            }}>
              Attestation (Type: I CONFIRM DELETE) <span style={{ color: 'var(--danger)' }}>*</span>
            </label>
            <Input
              value={attestation}
              onChangeText={setAttestation}
              placeholder="Type: I CONFIRM DELETE"
              required
              fullWidth
            />
            <p style={{ margin: '6px 0 0 0', fontSize: '12px', color: 'var(--text-secondary)' }}>
              Type exactly "<strong>I CONFIRM DELETE</strong>" to proceed
            </p>
          </div>

          {/* Action Buttons */}
          <div style={{
            display: 'flex',
            gap: '12px',
            paddingTop: '20px',
            borderTop: '1px solid var(--border-light)'
          }}>
            <Button
              type="submit"
              variant="danger"
              fullWidth
              disabled={loading || attestation !== "I CONFIRM DELETE"}
            >
              {loading ? 'Deleting...' : 'Confirm Delete'}
            </Button>
            <Button
              type="button"
              variant="secondary"
              onClick={handleClose}
              fullWidth
              disabled={loading}
            >
              Cancel
            </Button>
          </div>
        </div>
      </form>
    </Modal>
  );
}

export default DeleteConfirmationModal;
