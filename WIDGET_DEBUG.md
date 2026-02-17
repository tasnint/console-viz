# Debug: Why BarChart Shows as Table

## Possible Issues

### 1. **Widget Creation Failing Silently**

If the barchart widget fails to create (e.g., no numeric data in JSON), it might be falling back to table.

**Check:** Look for error messages when running:
```bash
go run ./cmd/console-viz/main.go "C:\Users\tanis\Downloads\MOCK_DATA.json" --widget=barchart
```

### 2. **Flag Not Being Parsed**

The `--widget` flag might not be parsed correctly.

**Test:** Try with explicit flag:
```bash
go run ./cmd/console-viz/main.go --file="C:\Users\tanis\Downloads\MOCK_DATA.json" --widget=barchart
```

### 3. **JSON Data Format**

The JSON might not have numeric arrays that barchart can use.

**Check your JSON format:**
- ✅ Works: `[10, 20, 30]` or `{"data": [10, 20, 30]}`
- ❌ Won't work: `{"name": "value"}` (no numbers)

### 4. **Default Widget**

If widget creation fails, it might default to table.

---

## ✅ **Fixed Issues**

1. ✅ Better error messages - now shows why widget creation failed
2. ✅ Exit on error - no silent failures
3. ✅ Clear error messages about data format

---

## 🔍 **Debug Steps**

1. **Check if flag is parsed:**
   ```bash
   # Should show barchart widget
   go run ./cmd/console-viz/main.go "file.json" --widget=barchart
   ```

2. **Check JSON format:**
   - Open your JSON file
   - Make sure it has numeric arrays: `[1, 2, 3]` or `{"data": [1, 2, 3]}`

3. **Try with CSV instead:**
   ```bash
   go run ./cmd/console-viz/main.go data/sample_sales.csv --widget=barchart
   ```

4. **Check error messages:**
   - The updated code now shows clear errors if widget creation fails
   - Look for: "Error: Failed to create widget 'barchart': ..."

---

## 🎯 **Quick Test**

```bash
# Test with sample CSV (should work)
go run ./cmd/console-viz/main.go data/sample_sales.csv --widget=barchart

# Test with your JSON
go run ./cmd/console-viz/main.go "C:\Users\tanis\Downloads\MOCK_DATA.json" --widget=barchart
```

If barchart fails, you'll now see a clear error message explaining why!
