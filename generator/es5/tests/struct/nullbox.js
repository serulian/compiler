$module('nullbox', function () {
  var $static = this;
  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'Value', {
        get: function () {
          if (this.$lazychecked['Value']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['Value'], $g.____testlib.basictypes.String, true, 'Value');
            this.$lazychecked['Value'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['Value'], $g.____testlib.basictypes.String);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['Value'] = this.Value;
        return $promise.resolve($t.box(mapped, $g.____testlib.basictypes.Mapping($t.any)));
      };
      return instance;
    };
    $instance.Mapping = function () {
      return $promise.resolve($t.box(this[BOXED_DATA_PROPERTY], $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['Value'], right[BOXED_DATA_PROPERTY]['Value'], $g.____testlib.basictypes.String));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'Value', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['Value'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['Value'] = value;
      },
    });
    this.$typesig = function () {
      return $t.createtypesig(['Value', 5, $g.____testlib.basictypes.String.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.nullbox.SomeStruct).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.nullbox.SomeStruct).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
  });

  $static.TEST = function () {
    var s;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.nullbox.SomeStruct.new().then(function ($result0) {
              $temp0 = $result0;
              $result = ($temp0, $temp0.Value = null, $temp0);
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            s = $result;
            $promise.resolve(s.Value == null).then(function ($result0) {
              return ($promise.shortcircuit($result0, true) || s.Mapping()).then(function ($result2) {
                return ($promise.shortcircuit($result0, true) || $result2.$index($t.box('Value', $g.____testlib.basictypes.String))).then(function ($result1) {
                  $result = $t.box($result0 && ($result1 == null), $g.____testlib.basictypes.Boolean);
                  $current = 2;
                  $continue($resolve, $reject);
                  return;
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
