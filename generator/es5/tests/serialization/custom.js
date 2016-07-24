$module('custom', function () {
  var $static = this;
  this.$class('CustomJSON', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.Get = function () {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $g.custom.CustomJSON.new().then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
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
    };
    $instance.Stringify = function (value) {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $g.____testlib.basictypes.JSON.Get().then(function ($result1) {
                return $result1.Stringify(value).then(function ($result0) {
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
    $instance.Parse = function (value) {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $g.____testlib.basictypes.JSON.Get().then(function ($result1) {
                return $result1.Parse(value).then(function ($result0) {
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
    this.$typesig = function () {
      return $t.createtypesig(['Get', 1, $g.____testlib.basictypes.Function($g.custom.CustomJSON).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Parse', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.custom.CustomJSON).$typeref()]);
    };
  });

  this.$struct('AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherBool) {
      var instance = new $static();
      instance.$unboxed = false;
      instance[BOXED_DATA_PROPERTY] = {
        AnotherBool: AnotherBool,
      };
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'AnotherBool', 'AnotherBool', function () {
      return $g.____testlib.basictypes.Boolean;
    }, true, function () {
      return $g.____testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.custom.AnotherStruct).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.custom.AnotherStruct).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.custom.AnotherStruct).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
  });

  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField, AnotherField, SomeInstance) {
      var instance = new $static();
      instance.$unboxed = false;
      instance[BOXED_DATA_PROPERTY] = {
        SomeField: SomeField,
        AnotherField: AnotherField,
        SomeInstance: SomeInstance,
      };
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'SomeField', function () {
      return $g.____testlib.basictypes.Integer;
    }, true, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'AnotherField', 'AnotherField', function () {
      return $g.____testlib.basictypes.Boolean;
    }, true, function () {
      return $g.____testlib.basictypes.Boolean;
    }, false);
    $t.defineStructField($static, 'SomeInstance', 'SomeInstance', function () {
      return $g.custom.AnotherStruct;
    }, true, function () {
      return $g.custom.AnotherStruct;
    }, false);
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.custom.SomeStruct).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.custom.SomeStruct).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.custom.SomeStruct).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
  });

  $static.TEST = function () {
    var jsonString;
    var parsed;
    var s;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.custom.AnotherStruct.new($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result1) {
              $temp0 = $result1;
              return $g.custom.SomeStruct.new($t.box(2, $g.____testlib.basictypes.Integer), $t.box(false, $g.____testlib.basictypes.Boolean), ($temp0, $temp0)).then(function ($result0) {
                $temp1 = $result0;
                $result = ($temp1, $temp1);
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
            s = $result;
            jsonString = $t.box('{"AnotherField":false,"SomeField":2,"SomeInstance":{"AnotherBool":true}}', $g.____testlib.basictypes.String);
            $g.custom.SomeStruct.Parse($g.custom.CustomJSON)(jsonString).then(function ($result0) {
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
            parsed = $result;
            $g.____testlib.basictypes.Integer.$equals(parsed.SomeField, $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result2) {
              return $promise.resolve($t.unbox($result2)).then(function ($result1) {
                return $promise.resolve($result1 && !$t.unbox(parsed.AnotherField)).then(function ($result0) {
                  $result = $t.box($result0 && $t.unbox(parsed.SomeInstance.AnotherBool), $g.____testlib.basictypes.Boolean);
                  $current = 3;
                  $continue($resolve, $reject);
                  return;
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
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
