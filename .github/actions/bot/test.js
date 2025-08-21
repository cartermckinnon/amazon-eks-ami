// Simple test file to validate bot command parsing
// Run with: node test.js

const bot = require('./index.js');

// Mock objects for testing
const mockCore = {
    setFailed: (error) => console.error('FAILED:', error)
};

const mockGithub = {
    rest: {
        issues: {
            createComment: (params) => {
                console.log('Would create comment:', params.body);
                return Promise.resolve();
            }
        }
    }
};

// Test cases
const testCases = [
    {
        name: 'Repro command with NodeConfig options',
        payload: {
            comment: {
                user: { login: 'testuser' },
                author_association: 'MEMBER',
                html_url: 'https://github.com/test/repo/issues/123#issuecomment-456',
                created_at: '2025-01-01T00:00:00Z',
                body: `/repro
+nodeconfig instance.localStorage=RAID0
+nodeconfig kubelet.maxPods=110`
            },
            issue: {
                number: 123,
                pull_request: null // This is an issue, not a PR
            },
            repository: {
                owner: { login: 'awslabs' },
                name: 'amazon-eks-ami'
            }
        }
    },
    {
        name: 'Repro command with AMI release tag and NodeConfig',
        payload: {
            comment: {
                user: { login: 'testuser' },
                author_association: 'MEMBER',
                html_url: 'https://github.com/test/repo/issues/2386#issuecomment-789',
                created_at: '2025-01-01T00:00:00Z',
                body: `/repro
+ami v20250620
+nodeconfig instance.localStorage=RAID0`
            },
            issue: {
                number: 2386,
                pull_request: null
            },
            repository: {
                owner: { login: 'awslabs' },
                name: 'amazon-eks-ami'
            }
        }
    },
    {
        name: 'CI command on PR',
        payload: {
            comment: {
                user: { login: 'testuser' },
                author_association: 'MEMBER',
                html_url: 'https://github.com/test/repo/pull/456#issuecomment-789',
                created_at: '2025-01-01T00:00:00Z',
                body: `/ci test
+os_distros al2023`
            },
            issue: {
                number: 456,
                pull_request: { url: 'https://api.github.com/repos/test/repo/pulls/456' }
            },
            repository: {
                owner: { login: 'awslabs' },
                name: 'amazon-eks-ami'
            }
        }
    },
    {
        name: 'Repro command on PR (should fail)',
        payload: {
            comment: {
                user: { login: 'testuser' },
                author_association: 'MEMBER',
                html_url: 'https://github.com/test/repo/pull/789#issuecomment-012',
                created_at: '2025-01-01T00:00:00Z',
                body: `/repro
+nodeconfig instance.localStorage=RAID0`
            },
            issue: {
                number: 789,
                pull_request: { url: 'https://api.github.com/repos/test/repo/pulls/789' }
            },
            repository: {
                owner: { login: 'awslabs' },
                name: 'amazon-eks-ami'
            }
        }
    },
    {
        name: 'CI command on issue (should fail)',
        payload: {
            comment: {
                user: { login: 'testuser' },
                author_association: 'MEMBER',
                html_url: 'https://github.com/test/repo/issues/321#issuecomment-654',
                created_at: '2025-01-01T00:00:00Z',
                body: `/ci test`
            },
            issue: {
                number: 321,
                pull_request: null
            },
            repository: {
                owner: { login: 'awslabs' },
                name: 'amazon-eks-ami'
            }
        }
    }
];

// Run tests
async function runTests() {
    console.log('Running bot tests...\n');
    
    for (const testCase of testCases) {
        console.log(`\n=== Test: ${testCase.name} ===`);
        
        const mockContext = {
            payload: testCase.payload,
            runId: '12345'
        };
        
        try {
            await bot(mockCore, mockGithub, mockContext, 'test-uuid-123');
            console.log('✅ Test completed successfully');
        } catch (error) {
            console.error('❌ Test failed:', error);
        }
    }
    
    console.log('\n=== All tests completed ===');
}

// Only run tests if this file is executed directly
if (require.main === module) {
    runTests();
}

module.exports = { runTests };
