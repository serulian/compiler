$module('asyncstruct', function () {
  var $static = this;
  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Foo, Bar) {
      var instance = new $static();
      instance.$unboxed = false;
      instance[BOXED_DATA_PROPERTY] = {
        Foo: Foo,
        Bar: Bar,
      };
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'Foo', 'Foo', function () {
      return $g.____testlib.basictypes.Integer;
    }, true, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'Bar', 'Bar', function () {
      return $g.____testlib.basictypes.Integer;
    }, true, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.asyncstruct.SomeStruct).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.asyncstruct.SomeStruct).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.asyncstruct.SomeStruct).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
  });

  $static.DoSomethingAsync = $t.workerwrap('dd7aa26e', function (s) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            s.Foo.String().then(function ($result1) {
              return s.Bar.String().then(function ($result2) {
                return $g.____testlib.basictypes.String.$plus($result1, $result2).then(function ($result0) {
                  $result = $result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
  $static.TEST = function () {
    var $result;
    var vle;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.asyncstruct.SomeStruct.new($t.box(1, $g.____testlib.basictypes.Integer), $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $promise.translate($g.asyncstruct.DoSomethingAsync($result1)).then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            vle = $result;
            $g.____testlib.basictypes.String.$equals(vle, $t.box("12", $g.____testlib.basictypes.String)).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
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
