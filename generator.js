// FSAgent JSON Generator

const FSAgentGenerator = {};

FSAgentGenerator.workspaceToCode = function(workspace) {
    const jobs = [];
    const topBlocks = workspace.getTopBlocks(false);
    
    for (let block of topBlocks) {
        if (block.type === 'fsagent_job') {
            const job = FSAgentGenerator.blockToJob(block);
            if (job) {
                jobs.push(job);
            }
        }
    }
    
    return JSON.stringify(jobs, null, 2);
};

FSAgentGenerator.blockToJob = function(block) {
    if (block.type !== 'fsagent_job') return null;
    
    const job = {
        folder: block.getFieldValue('FOLDER') || '/tmp',
        match: block.getFieldValue('MATCH') || '.*',
        verbose: block.getFieldValue('VERBOSE') === 'TRUE',
        debounce: block.getFieldValue('DEBOUNCE') === 'TRUE',
        onlynew: block.getFieldValue('ONLYNEW') === 'TRUE'
    };
    
    // Trigger verarbeiten
    const triggerBlock = block.getInputTargetBlock('TRIGGER');
    if (triggerBlock) {
        const trigger = FSAgentGenerator.blockToTrigger(triggerBlock);
        Object.assign(job, trigger);
    } else {
        job.trigger = 'fsevent';
    }
    
    // Actions verarbeiten
    const actionsBlock = block.getInputTargetBlock('ACTIONS');
    if (actionsBlock) {
        job.action = FSAgentGenerator.blockToActions(actionsBlock);
    } else {
        job.action = [];
    }
    
    return job;
};

FSAgentGenerator.blockToTrigger = function(block) {
    switch (block.type) {
        case 'trigger_fsevent':
            return { trigger: 'fsevent' };
        case 'trigger_ticker':
            return { 
                trigger: 'ticker',
                ticker: parseInt(block.getFieldValue('INTERVAL')) || 2000
            };
        case 'trigger_http':
            return { 
                trigger: 'http',
                port: block.getFieldValue('PORT') || '8080'
            };
        default:
            return { trigger: 'fsevent' };
    }
};

FSAgentGenerator.blockToActions = function(block) {
    const actions = [];
    let currentBlock = block;
    
    while (currentBlock) {
        const action = FSAgentGenerator.blockToAction(currentBlock);
        if (action) {
            actions.push(action);
        }
        currentBlock = currentBlock.getNextBlock();
    }
    
    return actions.length === 1 ? actions[0] : actions;
};

FSAgentGenerator.blockToAction = function(block) {
    let action = {};
    
    switch (block.type) {
        case 'action_sleep':
            action = {
                do: 'sleep',
                config: {
                    time: parseInt(block.getFieldValue('TIME')) || 1000
                }
            };
            break;
            
        case 'action_mail':
            const toField = block.getFieldValue('TO') || '';
            const ccField = block.getFieldValue('CC') || '';
            const bccField = block.getFieldValue('BCC') || '';
            
            action = {
                do: 'mail',
                config: {
                    name: 'mail',
                    subject: block.getFieldValue('SUBJECT') || 'FSAgent Notification',
                    body: block.getFieldValue('BODY') || 'File processed',
                    from: block.getFieldValue('FROM') || '',
                    to: toField.split(',').map(s => s.trim()).filter(s => s),
                    user: block.getFieldValue('USER') || '',
                    pass: block.getFieldValue('PASS') || '',
                    server: block.getFieldValue('SERVER') || '',
                    port: parseInt(block.getFieldValue('PORT')) || 587
                }
            };
            
            if (ccField) {
                action.config.cc = ccField.split(',').map(s => s.trim()).filter(s => s);
            }
            if (bccField) {
                action.config.bcc = bccField.split(',').map(s => s.trim()).filter(s => s);
            }
            break;
            
        case 'action_move':
            action = {
                do: 'move',
                config: {
                    name: block.getFieldValue('DESTINATION') || '/tmp/$file'
                }
            };
            break;
            
        case 'action_copy':
            action = {
                do: 'copy',
                config: {
                    name: block.getFieldValue('DESTINATION') || '/tmp/$file'
                }
            };
            break;
            
        case 'action_delete':
            action = {
                do: 'delete',
                config: {}
            };
            break;
            
        case 'action_checksize':
            action = {
                do: 'checksize',
                config: {
                    minSize: block.getFieldValue('MIN_SIZE') || '1KB',
                    maxSize: block.getFieldValue('MAX_SIZE') || '10MB'
                }
            };
            break;
            
        case 'action_http_post':
            const headers = {};
            const headersStr = block.getFieldValue('HEADERS') || '';
            if (headersStr) {
                headersStr.split(',').forEach(header => {
                    const [key, value] = header.split(':').map(s => s.trim());
                    if (key && value) {
                        headers[key] = value;
                    }
                });
            }
            
            action = {
                do: 'httppostrequest',
                config: {
                    url: block.getFieldValue('URL') || '',
                    headers: headers
                }
            };
            break;
            
        case 'action_openai':
            action = {
                do: 'openai',
                config: {
                    apiKey: block.getFieldValue('API_KEY') || '',
                    model: block.getFieldValue('MODEL') || 'gpt-3.5-turbo',
                    action: block.getFieldValue('ACTION') || 'analyze',
                    prompt: block.getFieldValue('PROMPT') || 'Analyze this file content',
                    outputFile: block.getFieldValue('OUTPUT_FILE') || ''
                }
            };
            break;
            
        case 'action_notification':
            action = {
                do: 'notification',
                config: {
                    webhookUrl: block.getFieldValue('WEBHOOK_URL') || '',
                    platform: block.getFieldValue('PLATFORM') || 'slack',
                    title: block.getFieldValue('TITLE') || 'FSAgent Notification',
                    message: block.getFieldValue('MESSAGE') || 'File processed: $file',
                    color: block.getFieldValue('COLOR') || 'good',
                    includeFile: block.getFieldValue('INCLUDE_FILE') === 'TRUE'
                }
            };
            break;
            
        case 'action_database':
            action = {
                do: 'database',
                config: {
                    driver: block.getFieldValue('DRIVER') || 'sqlite3',
                    connectionString: block.getFieldValue('CONNECTION') || './fsagent.db',
                    table: block.getFieldValue('TABLE') || 'fsagent_files',
                    action: block.getFieldValue('ACTION') || 'log',
                    createTable: true
                }
            };
            break;
            
        case 'action_fileprocessing':
            action = {
                do: 'fileprocessing',
                config: {
                    action: block.getFieldValue('ACTION') || 'hash',
                    hashType: block.getFieldValue('HASH_TYPE') || 'md5',
                    regexPattern: block.getFieldValue('REGEX_PATTERN') || '',
                    outputFile: block.getFieldValue('OUTPUT_FILE') || ''
                }
            };
            break;
            
        case 'action_ftp':
            action = {
                do: 'ftp',
                config: {
                    host: block.getFieldValue('HOST') || '',
                    port: parseInt(block.getFieldValue('PORT')) || 22,
                    protocol: block.getFieldValue('PROTOCOL') || 'sftp',
                    username: block.getFieldValue('USERNAME') || '',
                    password: block.getFieldValue('PASSWORD') || '',
                    remoteDir: block.getFieldValue('REMOTE_DIR') || '/upload/',
                    rename: block.getFieldValue('RENAME') || '',
                    keepFile: block.getFieldValue('KEEP_FILE') === 'TRUE'
                }
            };
            break;
            
        default:
            return null;
    }
    
    // onSuccess und onFailure Actions verarbeiten
    const onSuccessBlock = block.getInputTargetBlock('ON_SUCCESS');
    if (onSuccessBlock) {
        action.onSuccess = FSAgentGenerator.blockToActions(onSuccessBlock);
    }
    
    const onFailureBlock = block.getInputTargetBlock('ON_FAILURE');
    if (onFailureBlock) {
        action.onFailure = FSAgentGenerator.blockToActions(onFailureBlock);
    }
    
    return action;
};
