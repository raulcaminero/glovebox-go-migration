import React, { useState } from 'react';
import { X } from 'lucide-react';

const jwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyJ9.ah_KrNhrpocJD5ZKfyjoKMdCZMst4lvLSayT5B_d8Bk";

interface CreateModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreated: () => void;
}

export const CreateModal: React.FC<CreateModalProps> = ({ isOpen, onClose, onCreated }) => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [submitting, setSubmitting] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !email) {
      alert('Please fill in Name and Email');
      return;
    }
    setSubmitting(true);
    try {
      await fetch('/api/v1/policyholders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${jwtToken}`,
        },
        body: JSON.stringify({ full_name: name, email, phone }),
      });
      setName('');
      setEmail('');
      setPhone('');
      onCreated();
      onClose();
    } catch (err: any) {
      alert(`Failed to create policyholder: ${err.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="modal-overlay active">
      <div className="modal-card">
        <div className="modal-header">
          <h3>Create New Policyholder (Go API)</h3>
          <button className="btn btn-secondary text-xs" onClick={onClose}>
            <X size={14} />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="flex-col gap-3">
          <input
            type="text"
            className="input-field"
            placeholder="Full Name (e.g. John Doe)"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <input
            type="email"
            className="input-field"
            placeholder="Email (e.g. john@example.com)"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <input
            type="text"
            className="input-field"
            placeholder="Phone Number"
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
          />
          <div className="flex-row justify-end gap-2 mt-2">
            <button type="button" className="btn btn-secondary" onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className="btn" disabled={submitting}>
              {submitting ? 'Creating...' : 'Create Policyholder'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
