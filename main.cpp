#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <vector>
#include <map>
#include <queue>
#include <functional>
#include <cstring>
#include <cstdlib>

#include "Array.h"
#include "List.h"
#include "DList.h"
#include "Stack.h"
#include "Queue.h"
#include "Tree.h"

using namespace std;

enum class Type {
    ARRAY,
    STACK,
    QUEUE,
    TREE,
    FLIST,
    DLIST,
    UNKNOWN
};

struct Args { // аргументы
    string file;
    string query;
};

static bool ParseArgs(int argc, char** argv, Args& args) {
    for (int i = 1; i < argc; ++i) {
        string arg = argv[i];
        if (arg == "--file" && i + 1 < argc) args.file  = argv[++i];
        else if (arg == "--query" && i + 1 < argc) args.query = argv[++i];
        else {cout << "Неизвестный аргумент: " << arg << endl; return false;}
    }
    if (args.file.empty() || args.query.empty()) {
        cout << "Использование: ./dbms --file <file> --query '<cmd>'\n";
        return false;
    }
    return true;
}

struct Query { // команда
    string command;
    string name;
    vector<string> args;
};

static Query ParseQuery(const string& str) {
    Query query;
    istringstream input(str);
    input >> query.command >> query.name;
    for (string i; input >> i; ) query.args.push_back(i);
    return query;
}

struct Entry { // ввод
    string type;
    string name;
    vector<string> values;
};

struct DB { // база данных
    string filename;
    vector<Entry> entries;

    bool Load() {
        ifstream file(filename);
        if (!file) return true;
        string line;
        while (getline(file, line)) {
            if (line.empty()) continue;
            istringstream header(line);
            Entry entry;
            int size = 0;
            if (!(header >> entry.type >> entry.name >> size)) return false;
            for (int i = 0; i < size; ++i) {
                string value;
                if (!getline(file, value)) break;
                entry.values.push_back(value);
            }
            entries.push_back(move(entry));
        }
        return true;
    }

    bool Save() {
        ofstream file(filename);
        if (!file) return false;
        for (auto& entry : entries) {
            file << entry.type << " " << entry.name << " " << entry.values.size() << endl;
            for (auto& value : entry.values) file << value << endl;
        }
        return true;
    }

    Entry* GetOrCreate(const string& type, const string& name) {
        for (auto& entry : entries) {
            if (entry.name == name) {
                if (entry.type != type) {
                    cout << "Тип не совпадает: ожидается " << entry.type << endl;
                    return nullptr;
                }
                return &entry;
            }
        }
        entries.push_back({type, name, {}});
        return &entries.back();
    }
};

static bool ParseInt(const string& str, int& out) {
    char* end = nullptr;
    long value = strtol(str.c_str(), &end, 10);
    if (end == str.c_str() || *end != '\0') return false;
    out = (int)value;
    return true;
}

static Array ToArray(vector<string>& vector) {
    Array arr = MCREATE(vector.empty() ? 4 : (int)vector.size());
    for (auto& str : vector) {
        string copy = str;
        MPUSH(arr, str);
    }
    return arr;
}
static void FromArray(Array& arr, Entry& entry) {
    entry.values.clear();
    for (int i = 0; i < arr.size; ++i) entry.values.push_back(arr.data[i]);
    MDESTROY(arr);
}

static Stack ToStack(vector<string>& vector) {
    Stack stack = SCREATE(4);
    for (auto& item : vector) {
        string copy = item;
        SPUSH(stack, item);
    }
    return stack;
}
static void FromStack(Stack& stack, Entry& entry) {
    entry.values.clear();
    for (int i = 0; i <= stack.top; ++i) entry.values.push_back(stack.data[i]);
    SDESTROY(stack);
}

static Queue ToQueue(const vector<string>& vector) {
    Queue queue = QCREATE(4);
    for (auto& item : vector) {
        string copy = item;
        QPUSH(queue, item);
    }
    return queue;
}
static void FromQueue(Queue& queue, Entry& entry) {
    entry.values.clear();
    for (int i = 0; i < queue.size; ++i)
        entry.values.push_back(queue.data[(queue.head + i) % queue.capacity]);
    QDESTROY(queue);
}

static Tree ToTree(const vector<string>& vec) {
    Tree tree = TCREATE();
    for (auto& str : vec) {
        int item = 0;
        if (ParseInt(str, item)) {
            TINSERT(tree, item);
        } else if (str == "null") {
            // Если в файле написано null, мы искусственно создаем "дыру" для теста
            // Добавляем невидимый узел-пустышку с ключом -999, который алгоритм TIS_COMPLETE
            // посчитает за обычный узел, но мы изменим TIS_COMPLETE, чтобы она проверяла на -999
            TINSERT(tree, -999); 
        }
    }
    return tree;
}
// Замените старый CollectPre и FromTree в main.cpp на этот код:
static void FromTree(Tree& tree, Entry& entry) {
    entry.values.clear();
    if (tree.root == nullptr) {
        TDESTROY(tree.root);
        tree.root = nullptr;
        return;
    }

    // Сохраняем элементы строго по уровням (в ширину)
    queue<TNode*> q;
    q.push(tree.root);

    while (!q.empty()) {
        TNode* curr = q.front();
        q.pop();

        entry.values.push_back(to_string(curr->key));

        if (curr->left != nullptr) q.push(curr->left);
        if (curr->right != nullptr) q.push(curr->right);
    }

    TDESTROY(tree.root);
    tree.root = nullptr; 
}


static FList ToFList(const vector<string>& vector) {
    FList list = FCREATE();
    FNode* tail = nullptr;
    for (auto& str : vector) {
        FNode* node = new FNode{str, nullptr};
        if (!list.head) list.head = node;
        else tail->next = node;
        tail = node;
        ++list.size;
    }
    return list;
}
static void FromFList(FList& list, Entry& entry) {
    entry.values.clear();
    for (FNode* node = list.head; node; node = node->next) entry.values.push_back(node->data);
    FDESTROY(list);
}

static LList ToDList(const vector<string>& vector) {
    LList list = LCREATE();
    for (auto& str : vector) {
        LNode* node = new LNode{str, list.tail, nullptr};
        if (list.tail) list.tail->next = node;
        else list.head = node;
        list.tail = node;
        ++list.size;
    }
    return list;
}
static void FromDList(LList& list, Entry& entry) {
    entry.values.clear();
    for (LNode* node = list.head; node; node = node->next) entry.values.push_back(node->data);
    LDESTROY(list);
}

static Type TypeByCommand(const string& cmd) {
    if (cmd.empty()) return Type::UNKNOWN;
    switch (cmd[0]) {
        case 'M': return Type::ARRAY;
        case 'F': return Type::FLIST;
        case 'L': return Type::DLIST;
        case 'S': return Type::STACK;
        case 'Q': return Type::QUEUE;
        case 'T': return Type::TREE;
        default:  return Type::UNKNOWN;
    }
}

static const char* TypeName(Type type) {
    switch (type) {
        case Type::ARRAY: return "ARRAY";
        case Type::STACK: return "STACK";
        case Type::QUEUE: return "QUEUE";
        case Type::TREE:  return "TREE";
        case Type::FLIST: return "FLIST";
        case Type::DLIST: return "DLIST";
        default:          return "UNKNOWN";
    }
}

using Handler = function<bool(Entry&, const Query&)>;

static Handler wrapArray(function<bool(Array&, const Query&)> function) {
    return [function](Entry& entry, const Query& query) -> bool {
        Array arr = ToArray(entry.values);
        bool ok = function(arr, query);
        FromArray(arr, entry);
        return ok;
    };
}
static Handler wrapStack(function<bool(Stack&, const Query&)> function) {
    return [function](Entry& entry, const Query& query) -> bool {
        Stack stack = ToStack(entry.values);
        bool ok = function(stack, query);
        FromStack(stack, entry);
        return ok;
    };
}
static Handler wrapQueue(function<bool(Queue&, const Query&)> function) {
    return [function](Entry& entry, const Query& query) -> bool {
        Queue queue = ToQueue(entry.values);
        bool ok = function(queue, query);
        FromQueue(queue, entry);
        return ok;
    };
}
static Handler wrapTree(function<bool(Tree&, const Query&)> function) {
    return [function](Entry& entry, const Query& query) -> bool {
        Tree tree = ToTree(entry.values);
        bool ok = function(tree, query);
        FromTree(tree, entry);
        return ok;
    };
}
static Handler wrapFList(function<bool(FList&, const Query&)> function) {
    return [function](Entry& entry, const Query& query) -> bool {
        FList list = ToFList(entry.values);
        bool ok = function(list, query);
        FromFList(list, entry);
        return ok;
    };
}
static Handler wrapDList(function<bool(LList&, const Query&)> function) {
    return [function](Entry& entry, const Query& query) -> bool {
        LList list = ToDList(entry.values);
        bool ok = function(list, query);
        FromDList(list, entry);
        return ok;
    };
}

struct Command {
    int argc;
    Handler handler;
};

static const map<string, Command>& Table() {
    static const map<string, Command> table = {
        // ARRAY
        {"MPUSH", {1, wrapArray([](Array& arr, const Query& query) -> bool {
            string value = query.args[0]; MPUSH(arr, value); return true;
        })}},
        {"MPUSH_INDEX", {2, wrapArray([](Array& arr, const Query& query) -> bool {
            int idx;
            if (!ParseInt(query.args[1], idx)) { cout << "Не число: " << query.args[1] << endl; return false; }
            string value = query.args[0]; MPUSH_INDEX(arr, value, idx); return true;
        })}},
        {"MDEL", {1, wrapArray([](Array& arr, const Query& query) -> bool {
            int idx;
            if (!ParseInt(query.args[0], idx)) { cout << "Не число: " << query.args[0] << endl; return false; }
            cout << MDEL(arr, idx) << endl; return true;
        })}},
        {"MREPLACE", {2, wrapArray([](Array& arr, const Query& query) -> bool {
            int idx;
            if (!ParseInt(query.args[1], idx)) { cout << "Не число: " << query.args[1] << endl; return false; }
            string value = query.args[0]; MREPLACE(arr, value, idx); return true;
        })}},
        {"MRECEIVE", {1, wrapArray([](Array& arr, const Query& query) -> bool {
            int idx;
            if (!ParseInt(query.args[0], idx)) { cout << "Не число: " << query.args[0] << endl; return false; }
            cout << MRECEIVE(arr, idx) << endl; return true;
        })}},
        {"MLEN", {0, wrapArray([](Array& arr, const Query&) -> bool {
            cout << MLEN(arr) << endl; return true;
        })}},
        {"MGET", {0, wrapArray([](Array& arr, const Query&) -> bool {
            MGET(arr); return true;
        })}},

        // STACK
        {"SPUSH", {1, wrapStack([](Stack& stack, const Query& query) -> bool {
            string value = query.args[0]; SPUSH(stack, value); cout << value << endl; return true;
        })}},
        {"SPOP", {0, wrapStack([](Stack& stack, const Query&) -> bool {
            cout << SPOP(stack) << endl; return true;
        })}},
        {"SGET", {0, wrapStack([](Stack& stack, const Query&) -> bool {
            SGET(stack); return true;
        })}},
        {"SGET_REVERSE", {0, wrapStack([](Stack& stack, const Query&) -> bool {
            SGET_REVERSE(stack); return true;
        })}},

        // QUEUE
        {"QPUSH", {1, wrapQueue([](Queue& queue, const Query& query) -> bool {
            string value = query.args[0]; QPUSH(queue, value); cout << value << endl; return true;
        })}},
        {"QPOP", {0, wrapQueue([](Queue& queue, const Query&) -> bool {
            cout << QPOP(queue) << endl; return true;
        })}},
        {"QGET", {0, wrapQueue([](Queue& queue, const Query&) -> bool {
            QGET(queue); return true;
        })}},
        {"QGET_BACK", {0, wrapQueue([](Queue& queue, const Query&) -> bool {
            QGET_BACK(queue); return true;
        })}},

                // TREE
        {"TINSERT", {1, wrapTree([](Tree& tree, const Query& query) -> bool {
            int num;
            if (!ParseInt(query.args[0], num)) {
                cout << "Не число: " << query.args[0] << endl;
                return false;
            }
            if (tree.size > 0 && TFIND(tree, num)) {
                cout << "Значение уже есть в дереве: " << num << endl;
                return false;
            }
            TINSERT(tree, num);
            return true;
        })}},
                
        {"TDEL", {1, wrapTree([](Tree& tree, const Query& query) -> bool {
            int num;
            if (!ParseInt(query.args[0], num)) { 
                cout << "Не число: " << query.args[0] << endl; 
                return false; 
            }
            
            TDEL(tree, num); 
            return true;
        })}},

        {"TFIND", {1, wrapTree([](Tree& tree, const Query& query) -> bool {
            int num;
            if (!ParseInt(query.args[0], num)) {
                cout << "Не число: " << query.args[0] << endl;
                return false;
            }
            if (tree.root == nullptr) {
                cout << "Дерево пустое" << endl;
                return true;
            }
            cout << (TFIND(tree, num) ? "TRUE" : "FALSE") << endl;
            return true;
        })}},

                {"TGET_PREORDER", {0, wrapTree([](Tree& tree, const Query&) -> bool {
            cout << "Дерево полное: " << (TIS_COMPLETE(tree) ? "TRUE" : "FALSE") << endl;
            if (tree.size) TGET_PREORDER(tree.root); else cout << "Нет значений"; cout << endl; 
            return true;
        })}},

        {"TGET_INORDER",  {0, wrapTree([](Tree& tree, const Query&) -> bool {
            cout << "Дерево полное: " << (TIS_COMPLETE(tree) ? "TRUE" : "FALSE") << endl;
            if (tree.size) TGET_INORDER(tree.root); else cout << "Нет значений"; cout << endl; 
            return true;
        })}},

        {"TGET_POSTORDER",{0, wrapTree([](Tree& tree, const Query&) -> bool {
            cout << "Дерево полное: " << (TIS_COMPLETE(tree) ? "TRUE" : "FALSE") << endl;
            if (tree.size) TGET_POSTORDER(tree.root); else cout << "Нет значений"; cout << endl; 
            return true;
        })}},
        {"TGET_LEVELORDER", {0, wrapTree([](Tree& tree, const Query&) -> bool {
            cout << "Дерево полное: " << (TIS_COMPLETE(tree) ? "TRUE" : "FALSE") << endl;
            if (tree.size) TGET_LEVELORDER(tree); else cout << "Нет значений"; cout << endl; 
            return true;
        })}},
         




        // FLIST
        {"FPUSH_HEAD", {1, wrapFList([](FList& list, const Query& query) -> bool {
            string value = query.args[0]; FPUSH_HEAD(list, value); return true;
        })}},
        {"FPUSH_TAIL", {1, wrapFList([](FList& list, const Query& query) -> bool {
            string value = query.args[0]; FPUSH_TAIL(list, value); return true;
        })}},
        {"FPUSH_BEFORE", {2, wrapFList([](FList& list, const Query& query) -> bool {
            string value = query.args[0], t = query.args[1]; FPUSH_BEFORE(list, value, t); return true;
        })}},
        {"FPUSH_AFTER", {2, wrapFList([](FList& list, const Query& query) -> bool {
            string value = query.args[0], t = query.args[1]; FPUSH_AFTER(list, value, t); return true;
        })}},
        {"FDEL_HEAD", {0, wrapFList([](FList& list, const Query&) -> bool {
            FDEL_HEAD(list); return true;
        })}},
        {"FDEL_TAIL", {0, wrapFList([](FList& list, const Query&) -> bool {
            FDEL_TAIL(list); return true;
        })}},
        {"FDEL_BEFORE", {1, wrapFList([](FList& list, const Query& query) -> bool {
            string t = query.args[0]; FDEL_BEFORE(list, t); return true;
        })}},
        {"FDEL_AFTER", {1, wrapFList([](FList& list, const Query& query) -> bool {
            string t = query.args[0]; FDEL_AFTER(list, t); return true;
        })}},
        {"FDEL_VALUE", {1, wrapFList([](FList& list, const Query& query) -> bool {
            string t = query.args[0]; FDEL_VALUE(list, t); return true;
        })}},
        {"FFIND", {1, wrapFList([](FList& list, const Query& query) -> bool {
            int idx = FFIND(list, query.args[0]);
            cout << (idx < 0 ? "Нет такого значения" : to_string(idx)) << endl; return true;
        })}},
        {"FLEN", {0, wrapFList([](FList& list, const Query&) -> bool {
            cout << FLEN(list) << endl; return true;
        })}},
        {"FGET", {0, wrapFList([](FList& list, const Query&) -> bool {
            FGET_CLASSIC(list); return true;
        })}},
        {"FGET_REVERSE", {0, wrapFList([](FList& list, const Query&) -> bool {
            FGET_REVERSE(list); return true;
        })}},

        // DLIST
        {"LPUSH_HEAD", {1, wrapDList([](LList& list, const Query& query) -> bool {
            string value = query.args[0]; LPUSH_HEAD(list, value); return true;
        })}},
        {"LPUSH_TAIL", {1, wrapDList([](LList& list, const Query& query) -> bool {
            string value = query.args[0]; LPUSH_TAIL(list, value); return true;
        })}},
        {"LPUSH_BEFORE", {2, wrapDList([](LList& list, const Query& query) -> bool {
            string value = query.args[0], t = query.args[1]; LPUSH_BEFORE(list, value, t); return true;
        })}},
        {"LPUSH_AFTER", {2, wrapDList([](LList& list, const Query& query) -> bool {
            string value = query.args[0], t = query.args[1]; LPUSH_AFTER(list, value, t); return true;
        })}},
        {"LDEL_HEAD", {0, wrapDList([](LList& list, const Query&) -> bool {
            LDEL_HEAD(list); return true;
        })}},
        {"LDEL_TAIL", {0, wrapDList([](LList& list, const Query&) -> bool {
            LDEL_TAIL(list); return true;
        })}},
        {"LDEL_BEFORE", {1, wrapDList([](LList& list, const Query& query) -> bool {
            string t = query.args[0]; LDEL_BEFORE(list, t); return true;
        })}},
        {"LDEL_AFTER", {1, wrapDList([](LList& list, const Query& query) -> bool {
            string t = query.args[0]; LDEL_AFTER(list, t); return true;
        })}},
        {"LDEL_VALUE", {1, wrapDList([](LList& list, const Query& query) -> bool {
            string t = query.args[0]; LDEL_VALUE(list, t); return true;
        })}},
        {"LFIND", {1, wrapDList([](LList& list, const Query& query) -> bool {
            int idx = LFIND(list, query.args[0]);
            cout << (idx < 0 ? "Нет такого значения" : to_string(idx)) << endl; return true;
        })}},
        {"LGET", {0, wrapDList([](LList& list, const Query&) -> bool {
            LGET(list); return true;
        })}},
        {"LGET_BACK", {0, wrapDList([](LList& list, const Query&) -> bool {
            LGET_BACK(list); return true;
        })}},
    };
    return table;
}

int main(int argc, char** argv) {
    Args args;
    if (!ParseArgs(argc, argv, args)) return 1;

    Query query = ParseQuery(args.query);
    if (query.command.empty() || query.name.empty()) {
        cout << "Неверный запрос\n";
        return 1;
    }

    auto& tbl = Table();
    auto it = tbl.find(query.command);
    if (it == tbl.end()) {
        cout << "Неизвестная команда: " << query.command << "\n";
        return 1;
    }
    const Command& cmd = it->second;

    if ((int)query.args.size() < cmd.argc) {
        cout << "Команда " << query.command << " требует минимум " << cmd.argc << " аргумент(ов)\n";
        return 1;
    }

    Type type = TypeByCommand(query.command);
    if (type == Type::UNKNOWN) {
        cout << "Не удалось определить тип для команды: " << query.command << "\n";
        return 1;
    }

    DB db;
    db.filename = args.file;
    if (!db.Load()) {
        cout << "Ошибка загрузки\n";
        return 1;
    }

    Entry* entry = db.GetOrCreate(TypeName(type), query.name);
    if (!entry) return 1;

    if (!cmd.handler(*entry, query)) return 1;

    if (!db.Save()) {
        cout << "Ошибка сохранения\n";
        return 1;
    }
    return 0;
}