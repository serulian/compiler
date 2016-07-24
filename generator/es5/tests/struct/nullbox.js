$module('nullbox', function () {
  var $static = this;
  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.$unboxed = false;
      instance[BOXED_DATA_PROPERTY] = {
      };
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'Value', 'Value', function () {
      return $g.____testlib.basictypes.String;
    }, true, function () {
      return $g.____testlib.basictypes.String;
    }, true);
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.nullbox.SomeStruct).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.nullbox.SomeStruct).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.nullbox.SomeStruct).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
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
